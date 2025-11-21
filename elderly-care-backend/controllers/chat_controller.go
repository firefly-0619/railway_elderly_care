package controllers

import (
	"context"
	"elderly-care-backend/common/constants"
	. "elderly-care-backend/global"
	"elderly-care-backend/models"
	"elderly-care-backend/utils"
	"elderly-care-backend/vo"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 在生产环境中应该验证来源
	},
}

// 客户端连接信息
type Client struct {
	Conn      *websocket.Conn
	AccountID uint
}

// 连接管理器
type ChatManager struct {
	clients    map[uint]*Client
	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex //保证并发安全
}

// 创建连接管理器
func NewConnectionManager() *ChatManager {
	return &ChatManager{
		clients:    make(map[uint]*Client),
		broadcast:  make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// 启动连接管理器
func (manager *ChatManager) Start() {
	writer := KafkaOperators[constants.MESSAGE_TOPIC].Writer
	go manager.HandleMessage()
	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clients[client.AccountID] = client
			manager.mutex.Unlock()
			log.Printf("用户 %d 已连接", client.AccountID)

		case client := <-manager.unregister:
			manager.mutex.Lock()
			if _, exists := manager.clients[client.AccountID]; exists {
				delete(manager.clients, client.AccountID)
				client.Conn.Close()
				log.Printf("用户 %d 已断开连接", client.AccountID)
			}
			manager.mutex.Unlock()

		case message := <-manager.broadcast:
			// 给message两端的用户都发送，因为message经过了封装
			manager.SendMessage(message.From, message)
			manager.SendMessage(message.To, message)
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("JSON 序列化失败: %v", err)
				continue
			}
			//kafka写入消息
			_ = writer.WriteMessages(context.Background(), kafka.Message{
				Value: data,
			})
		}
	}
}

// 异步处理kafka消息，包括更新联系列表和存储
func (manager *ChatManager) HandleMessage() {
	reader := KafkaOperators[constants.MESSAGE_TOPIC].Reader
	for {
		kafkaMsg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka 读取消息失败: %v", err)
			continue
		}
		msg := models.Message{}
		err = json.Unmarshal(kafkaMsg.Value, &msg)
		if err != nil {
			log.Printf("JSON 反序列化失败: %v", err)
			continue
		}
		Db.Create(&msg)
		err = Db.Transaction(func(tx *gorm.DB) error {
			if err1 := manager.handleContactList(msg.From, msg.To, msg.Time, tx); err1 != nil {
				return err1
			}
			if err2 := manager.handleContactList(msg.To, msg.From, msg.Time, tx); err2 != nil {
				return err2
			}
			return nil
		})

		if err != nil {
			Logger.Error("Update ContactList error", zap.Error(err))
		}

		err = reader.CommitMessages(context.Background(), kafkaMsg)
		if err != nil {
			Logger.Error("Kafka 提交消息失败", zap.Error(err))
		}
	}
}

func (manager *ChatManager) handleContactList(accountID, contactID uint, lastChatTime time.Time, db *gorm.DB) error {
	if err := db.Create(&models.ContactList{
		AccountID:    accountID,
		ContactID:    contactID,
		LastChatTime: lastChatTime,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if err = db.Model(&models.ContactList{}).Where("account_id = ? and contact_id = ?", accountID, contactID).Update("last_chat_time", lastChatTime).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (manager *ChatManager) SendMessage(toID uint, message models.Message) {
	manager.mutex.RLock()
	// 发送给指定用户
	if targetClient, exists := manager.clients[toID]; exists {
		err := targetClient.Conn.WriteJSON(message)
		if err != nil {
			log.Printf("发送消息失败: %v", err)
			// 删除该用户，先释放读锁再加写锁
			manager.mutex.RUnlock()
			manager.mutex.Lock()
			delete(manager.clients, toID)
			targetClient.Conn.Close()
			log.Printf("用户 %s 已断开连接", toID)
			manager.mutex.Unlock()
		} else {
			manager.mutex.RUnlock()
		}
	} else {
		manager.mutex.RUnlock()
	}
}

// 处理 WebSocket 连接
func (manager *ChatManager) HandleWebSocket(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket 升级失败"})
		return
	}

	client := &Client{
		Conn:      conn,
		AccountID: utils.GetAccountIdInContext(c),
	}

	manager.register <- client

	defer func() {
		manager.unregister <- client
	}()

	// 读取消息
	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			break
		}
		ID, err := RedisClient.Incr(context.Background(), constants.CHAT_ID_COUNT).Result()
		if err != nil {
			log.Printf("获取消息ID失败: %v", err)
			break
		}

		// 设置发送者
		message.From = utils.GetAccountIdInContext(c)
		message.Time = time.Now()
		message.ID = uint(ID)
		// 广播消息
		manager.broadcast <- message
	}
}

// @Tags 聊天模块
// @Summary 聊天记录
// @Description 聊天记录
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param  accountID query uint true "聊天对象用户ID(用于上滑加载，避免一次性加载过多，默认返回最近的几条聊天信息)"
// @Param  messageID query uint false "消息ID"
// @Param  size query int false "返回数量"
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 500 {object} vo.ResponseVO "失败"
// @Router /chat/record [get]
func (manager *ChatManager) GetChatRecord(c *gin.Context) {
	accountIDStr := c.Query("accountID")
	messageIDStr := c.Query("messageID")
	size := c.Query("size")
	var accountID uint
	if accountIDStr != "" {
		t, _ := strconv.Atoi(accountIDStr)
		accountID = uint(t)
	}
	query := Db
	if messageIDStr != "" {
		t, _ := strconv.Atoi(messageIDStr)
		query = Db.Where("id < ?", t)
	}
	if size != "" {
		t, _ := strconv.Atoi(size)
		query = query.Limit(t)
	} else {
		query = query.Limit(constants.DEFAULT_PAGE_SIZE)
	}
	msgs := make([]models.Message, 0)
	if query.Where(models.Message{From: utils.GetAccountIdInContext(c), To: accountID}).
		Or(models.Message{From: accountID, To: utils.GetAccountIdInContext(c)}).Order("time desc").Find(&msgs).Error != nil {
		c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
	}
	c.JSON(http.StatusOK, vo.Success(msgs))
}

// @Tags 聊天模块
// @Summary 最近的联系人
// @Description 最近的联系人
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param lastChatTime query string false "聊天时间:格式，2025-11-09T13:49:06.732+08:00(用于上滑加载，避免一次性加载过多，默认返回最近的聊过天的几个联系人)"
// @Param size query int false "返回数量"
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 500 {object} vo.ResponseVO "失败"
// @Router /chat/contactList [get]
func (manager *ChatManager) GetRecentlyChatList(c *gin.Context) {
	size := utils.CoverStr2Int(c.Query("size"), 10)
	timeStr := c.Query("lastChatTime")
	contactList := make([]vo.ContactListVo, 0)
	query := Db.Model(&models.ContactList{}).
		Select("contact_list.contact_id,contact_list.last_chat_time, account.nickname, account.avatar").
		Joins(" join account  on account.id = contact_list.contact_id").Where("account_id = ?", utils.GetAccountIdInContext(c)).Order("last_chat_time desc").Limit(size)
	if timeStr != "" {
		lastChatTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			c.JSON(http.StatusBadGateway, vo.Fail(constants.PARAM_ERROR))
			return
		}
		query = query.Where("last_chat_time < ?", lastChatTime)
	}
	if query.Find(&contactList).Error != nil {
		c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
		return
	}
	c.JSON(http.StatusOK, vo.Success(contactList))

}
