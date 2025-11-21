package controllers

import (
	"elderly-care-backend/common/constants"
	"elderly-care-backend/dto"
	"elderly-care-backend/global"
	"elderly-care-backend/models"
	"elderly-care-backend/services"
	"elderly-care-backend/utils"
	"elderly-care-backend/vo"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SOSController struct {
	matchingService *services.TaskMatchingService
}

func NewSOSController() *SOSController {
	return &SOSController{
		matchingService: &services.TaskMatchingService{},
	}
}

// @Tags SOS模块
// @Summary 触发紧急求助
// @Description 用户触发SOS紧急求助
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.SOSRequest true "SOS求助请求参数"
// @Success 200 {object} vo.ResponseVO{data=vo.SOSResponseVO}
// @Router /sos/emergency [post]
func (sc *SOSController) TriggerEmergency(c *gin.Context) {
	var req dto.SOSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.Fail(constants.PARAM_ERROR))
		return
	}

	// 开始事务
	tx := global.Db.Begin()

	// 创建紧急任务
	task := models.Task{
		CreatorID:   req.UserID,
		Title:       "紧急求助",
		Description: req.Description,
		Category:    "emergency",
		Status:      "pending",
		Reward:      0,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Address:     req.Address,
	}

	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, vo.Fail(constants.SOS_CREATE_FAILED))
		return
	}

	// 创建SOS记录
	timeoutAt := time.Now().Add(5 * time.Minute)
	sosRecord := models.SOSRecord{
		UserID:      req.UserID,
		TaskID:      task.ID,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Address:     req.Address,
		Description: req.Description,
		Severity:    req.Severity,
		Status:      "pending",
		TimeoutAt:   &timeoutAt,
	}

	if err := tx.Create(&sosRecord).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, vo.Fail(constants.SOS_CREATE_FAILED))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, vo.Fail(constants.SOS_CREATE_FAILED))
		return
	}

	// 执行紧急匹配
	matches, _ := sc.matchingService.MatchEmergencyVolunteers(task.ID, req.Latitude, req.Longitude)

	c.JSON(http.StatusOK, vo.Success(vo.SOSResponseVO{
		SOSID:   sosRecord.ID,
		TaskID:  task.ID,
		Matches: matches,
		Timeout: 300, // 5分钟
	}))
}

// @Tags SOS模块
// @Summary 接受SOS求助
// @Description 志愿者接受SOS紧急求助
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sosId path int true "SOS记录ID"
// @Param request body dto.AcceptSOSRequest true "接受SOS请求参数"
// @Success 200 {object} vo.ResponseVO
// @Router /sos/{sosId}/accept [post]
func (sc *SOSController) AcceptSOS(c *gin.Context) {
	sosIDStr := c.Param("sosId")
	sosID, _ := strconv.ParseUint(sosIDStr, 10, 32)

	var req dto.AcceptSOSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.Fail(constants.PARAM_ERROR))
		return
	}

	// 获取SOS记录
	var sosRecord models.SOSRecord
	if err := global.Db.First(&sosRecord, uint(sosID)).Error; err != nil {
		c.JSON(http.StatusOK, vo.Fail(constants.SOS_NOT_EXIST))
		return
	}

	if sosRecord.Status != "pending" {
		c.JSON(http.StatusOK, vo.Fail(constants.SOS_ALREADY_RESOLVED))
		return
	}

	// 开始事务
	tx := global.Db.Begin()

	// 更新SOS记录状态
	if err := tx.Model(&sosRecord).Updates(map[string]interface{}{
		"status": "accepted",
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	// 更新关联的任务状态
	if err := tx.Model(&models.Task{}).Where("id = ?", sosRecord.TaskID).Updates(map[string]interface{}{
		"assignee_id": req.VolunteerID,
		"status":      "assigned",
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	c.JSON(http.StatusOK, vo.Success(nil))
}

// @Tags SOS模块
// @Summary 解决SOS求助
// @Description 标记SOS紧急求助为已解决
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sosId path int true "SOS记录ID"
// @Param request body dto.ResolveSOSRequest true "解决SOS请求参数"
// @Success 200 {object} vo.ResponseVO
// @Router /sos/{sosId}/resolve [put]
func (sc *SOSController) ResolveSOS(c *gin.Context) {
	sosIDStr := c.Param("sosId")
	sosID, _ := strconv.ParseUint(sosIDStr, 10, 32)

	var req dto.ResolveSOSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.Fail(constants.PARAM_ERROR))
		return
	}

	// 获取SOS记录
	var sosRecord models.SOSRecord
	if err := global.Db.First(&sosRecord, uint(sosID)).Error; err != nil {
		c.JSON(http.StatusOK, vo.Fail(constants.SOS_NOT_EXIST))
		return
	}

	if sosRecord.Status == "resolved" {
		c.JSON(http.StatusOK, vo.Fail(constants.SOS_ALREADY_RESOLVED))
		return
	}

	// 开始事务
	tx := global.Db.Begin()

	now := time.Now()
	// 更新SOS记录状态
	if err := tx.Model(&sosRecord).Updates(map[string]interface{}{
		"status":      "resolved",
		"resolved_at": &now,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	// 更新关联的任务状态
	if err := tx.Model(&models.Task{}).Where("id = ?", sosRecord.TaskID).Updates(map[string]interface{}{
		"status": "completed",
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	c.JSON(http.StatusOK, vo.Success(nil))
}

// @Tags SOS模块
// @Summary 获取当前SOS状态
// @Description 获取用户当前的SOS求助状态
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} vo.ResponseVO{data=vo.SOSRecordVO}
// @Router /sos/current [get]
func (sc *SOSController) GetCurrentSOS(c *gin.Context) {
	// 获取当前用户ID
	userID := utils.GetAccountIdInContext(c)

	var sosRecord models.SOSRecord
	if err := global.Db.Where("user_id = ? AND status IN (?)", userID, []string{"pending", "matching", "accepted", "in_progress"}).
		Order("created_at DESC").
		First(&sosRecord).Error; err != nil {
		// 没有进行中的SOS是正常的
		c.JSON(http.StatusOK, vo.Success(nil))
		return
	}

	sosVO := vo.SOSRecordVO{
		ID:          sosRecord.ID,
		UserID:      sosRecord.UserID,
		TaskID:      sosRecord.TaskID,
		Latitude:    sosRecord.Latitude,
		Longitude:   sosRecord.Longitude,
		Address:     sosRecord.Address,
		Description: sosRecord.Description,
		Severity:    sosRecord.Severity,
		Status:      sosRecord.Status,
		CreatedAt:   sosRecord.CreatedAt,
	}

	c.JSON(http.StatusOK, vo.Success(sosVO))
}
