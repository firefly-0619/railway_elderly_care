package controllers

import (
	"elderly-care-backend/common/constants"
	"elderly-care-backend/dto"
	"elderly-care-backend/global"
	"elderly-care-backend/models"
	"elderly-care-backend/services"
	"elderly-care-backend/vo"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	matchingService *services.TaskMatchingService
}

func NewTaskController() *TaskController {
	return &TaskController{
		matchingService: &services.TaskMatchingService{},
	}
}

// @Tags 任务模块
// @Summary 创建任务
// @Description 创建新的求助任务
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CreateTaskRequest true "创建任务请求参数"
// @Success 200 {object} vo.ResponseVO{data=vo.TaskMatchVO}
// @Router /tasks [post]
func (tc *TaskController) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.Fail(constants.PARAM_ERROR))
		return
	}

	// 创建任务 - 修复CreatorID类型
	task := models.Task{
		CreatorID:   uint(req.CreatorID), // 转换为uint
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Reward:      req.Reward,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Address:     req.Address,
	}

	if err := global.Db.Create(&task).Error; err != nil {
		c.JSON(http.StatusOK, vo.Fail(constants.TASK_CREATE_FAILED))
		return
	}

	// 执行匹配 - 这里task.ID已经是uint类型
	matches, _ := tc.matchingService.MatchVolunteersForTask(task.ID)

	// 返回VO对象
	taskVO := vo.TaskVO{
		ID:          task.ID,
		CreatorID:   task.CreatorID,
		Title:       task.Title,
		Description: task.Description,
		Category:    task.Category,
		Status:      task.Status,
		Reward:      task.Reward,
		Latitude:    task.Latitude,
		Longitude:   task.Longitude,
		Address:     task.Address,
		CreatedAt:   task.CreatedAt,
	}

	c.JSON(http.StatusOK, vo.Success(vo.TaskMatchVO{
		Task:    taskVO,
		Matches: matches,
	}))
}

// @Tags 任务模块
// @Summary 获取附近任务
// @Description 根据位置获取附近的任务列表
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param lat query number true "纬度"
// @Param lng query number true "经度"
// @Param radius query number false "搜索半径(米)" default(5000)
// @Param category query string false "任务分类"
// @Success 200 {object} vo.ResponseVO{data=[]vo.TaskVO}
// @Router /tasks/nearby [get]
func (tc *TaskController) GetNearbyTasks(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.Query("radius")
	category := c.Query("category")

	if latStr == "" || lngStr == "" {
		c.JSON(http.StatusOK, vo.Fail(constants.LOCATION_REQUIRED))
		return
	}

	lat, _ := strconv.ParseFloat(latStr, 64)
	lng, _ := strconv.ParseFloat(lngStr, 64)
	radius, _ := strconv.ParseFloat(radiusStr, 64)
	if radius == 0 {
		radius = 5000
	}

	// 修复：使用MySQL正确的距离计算公式
	query := `
        SELECT 
            t.*,
            (6371000 * ACOS(
                COS(? * PI() / 180) * COS(t.latitude * PI() / 180) * 
                COS((t.longitude - ?) * PI() / 180) + 
                SIN(? * PI() / 180) * SIN(t.latitude * PI() / 180)
            )) as distance,
            CASE WHEN t.category = 'emergency' THEN 1 ELSE 0 END as is_emergency
        FROM task t
        WHERE t.status = 'pending'
          AND (6371000 * ACOS(
                COS(? * PI() / 180) * COS(t.latitude * PI() / 180) * 
                COS((t.longitude - ?) * PI() / 180) + 
                SIN(? * PI() / 180) * SIN(t.latitude * PI() / 180)
            )) < ?`

	params := []interface{}{lat, lng, lat, lat, lng, lat, radius}
	if category != "" && category != "all" {
		query += " AND t.category = ?"
		params = append(params, category)
	}

	query += " ORDER BY is_emergency DESC, distance ASC LIMIT 50"

	var tasks []vo.TaskVO
	if err := global.Db.Raw(query, params...).Scan(&tasks).Error; err != nil {
		fmt.Printf("SQL错误: %v\n", err) // 添加错误日志
		c.JSON(http.StatusOK, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	c.JSON(http.StatusOK, vo.Success(tasks))
}

// @Tags 任务模块
// @Summary 接受任务
// @Description 志愿者接受任务
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param taskId path int true "任务ID"
// @Param request body dto.AcceptTaskRequest true "接受任务请求参数"
// @Success 200 {object} vo.ResponseVO
// @Router /tasks/{taskId}/accept [post]
func (tc *TaskController) AcceptTask(c *gin.Context) {
	taskIDStr := c.Param("taskId")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, vo.Fail("任务ID格式错误"))
		return
	}

	var req dto.AcceptTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("参数绑定错误: %v\n", err)
		c.JSON(http.StatusOK, vo.Fail(constants.PARAM_ERROR))
		return
	}

	// 1. 从 JWT 中获取用户ID（安全增强）
	userID, exists := c.Get("userID")
	if exists {
		// 如果 JWT 中有用户ID，使用 JWT 中的ID（更安全）
		jwtUserID := userID.(uint)
		fmt.Printf("从JWT获取用户ID: %d\n", jwtUserID)
		// 这里可以选择使用 JWT 中的用户ID，或者验证与请求体中的ID是否一致
	}

	// 2. 检查志愿者是否存在（数据验证）
	var volunteer models.Account
	if err := global.Db.First(&volunteer, req.VolunteerID).Error; err != nil {
		c.JSON(http.StatusOK, vo.Fail("志愿者用户不存在"))
		return
	}

	// 获取任务详情
	var task models.Task
	if err := global.Db.First(&task, uint(taskID)).Error; err != nil {
		fmt.Printf("任务查询错误: %v\n", err)
		c.JSON(http.StatusOK, vo.Fail(constants.TASK_NOT_EXIST))
		return
	}

	fmt.Printf("找到任务: ID=%d, Status=%s, CreatorID=%d\n", task.ID, task.Status, task.CreatorID)

	if task.Status != "pending" {
		fmt.Printf("任务状态不是pending: %s\n", task.Status)
		c.JSON(http.StatusOK, vo.Fail(constants.TASK_ALREADY_ACCEPTED))
		return
	}

	// 2. 检查志愿者不能接受自己的任务（数据验证）
	if task.CreatorID == req.VolunteerID {
		c.JSON(http.StatusOK, vo.Fail("不能接受自己的任务"))
		return
	}

	// 更新任务状态
	fmt.Printf("准备更新任务: assignee_id=%d, status=assigned\n", req.VolunteerID)
	if err := global.Db.Model(&task).Updates(map[string]interface{}{
		"assignee_id": uint(req.VolunteerID),
		"status":      "assigned",
	}).Error; err != nil {
		// 1. 详细的错误处理
		if strings.Contains(err.Error(), "foreign key constraint") {
			fmt.Printf("外键约束错误: %v\n", err)
			c.JSON(http.StatusOK, vo.Fail("志愿者用户不存在"))
			return
		}
		fmt.Printf("数据库更新错误: %v\n", err)
		c.JSON(http.StatusOK, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	fmt.Printf("任务接受成功! 任务ID=%d 被志愿者ID=%d 接受\n", task.ID, req.VolunteerID)
	c.JSON(http.StatusOK, vo.Success(nil))
}
