package services

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境需要严格检查
	},
}

type WebSocketService struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.Mutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// 启动WebSocket服务
func (ws *WebSocketService) Start() {
	for {
		select {
		case client := <-ws.register:
			ws.mutex.Lock()
			ws.clients[client] = true
			ws.mutex.Unlock()
			log.Println("客户端连接")

		case client := <-ws.unregister:
			ws.mutex.Lock()
			if _, ok := ws.clients[client]; ok {
				delete(ws.clients, client)
				client.Close()
			}
			ws.mutex.Unlock()
			log.Println("客户端断开")

		case message := <-ws.broadcast:
			ws.mutex.Lock()
			for client := range ws.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					client.Close()
					delete(ws.clients, client)
				}
			}
			ws.mutex.Unlock()
		}
	}
}

// 处理WebSocket连接
func (ws *WebSocketService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket升级失败:", err)
		return
	}

	ws.register <- conn

	defer func() {
		ws.unregister <- conn
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// 处理客户端消息（如位置更新）
		if messageType == websocket.TextMessage {
			ws.handleClientMessage(conn, p)
		}
	}
}

// 处理客户端消息
func (ws *WebSocketService) handleClientMessage(conn *websocket.Conn, message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Println("消息解析失败:", err)
		return
	}

	// 根据消息类型处理
	switch msg["type"] {
	case "location_update":
		ws.handleLocationUpdate(msg)
	case "navigation_request":
		ws.handleNavigationRequest(conn, msg)
	}
}

// 处理位置更新
func (ws *WebSocketService) handleLocationUpdate(msg map[string]interface{}) {
	// 这里可以保存位置到数据库，并广播给相关用户
	log.Printf("收到位置更新: %v", msg)

	// 广播给需要知道这个位置的其他用户
	broadcastMsg, _ := json.Marshal(map[string]interface{}{
		"type": "location_updated",
		"data": msg["data"],
	})
	ws.broadcast <- broadcastMsg
}

// 处理导航请求
func (ws *WebSocketService) handleNavigationRequest(conn *websocket.Conn, msg map[string]interface{}) {
	data := msg["data"].(map[string]interface{})
	startLat := data["start_lat"].(float64)
	startLng := data["start_lng"].(float64)
	endLat := data["end_lat"].(float64)
	endLng := data["end_lng"].(float64)

	// 计算导航信息
	locationService := &RealtimeLocationService{}
	navigation, err := locationService.CalculateRealTimeNavigation(startLat, startLng, endLat, endLng)

	response := map[string]interface{}{
		"type": "navigation_response",
		"data": navigation,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	responseMsg, _ := json.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, responseMsg)
}

// 广播消息给所有客户端
func (ws *WebSocketService) BroadcastMessage(messageType string, data interface{}) {
	msg := map[string]interface{}{
		"type": messageType,
		"data": data,
	}

	msgBytes, _ := json.Marshal(msg)
	ws.broadcast <- msgBytes
}
