package controllers

import (
	"elderly-care-backend/common/constants"
	. "elderly-care-backend/global"
	"elderly-care-backend/models"
	"elderly-care-backend/vo"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EvaluationController struct{}

// @Tags 评价模块
// @Summary 评价账号
// @Description 任务完成后对任务进行评价
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param  accountID query uint true "评价的账号ID"
// @Param  score query uint false "得分（1-5分）"
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 500 {object} vo.ResponseVO "失败"
// @Router /evaluation/account [put]
func (ce *EvaluationController) EvaluationAccount(c *gin.Context) {
	assignIDStr := c.Query("accountID")
	scoreStr := c.Query("score")
	if assignIDStr == "" || scoreStr == "" {
		c.JSON(http.StatusBadRequest, vo.Fail(constants.PARAM_ERROR))
	} else {
		assignID, err := strconv.Atoi(assignIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, vo.Fail(constants.PARAM_ERROR))
			return
		}
		score, err := strconv.Atoi(scoreStr)
		//判断得分是否正常
		if score < 1 || score > 5 {
			c.JSON(http.StatusBadRequest, vo.Fail(constants.EVALUATION_ERROR))
			return
		}
		//更新评价
		if Db.Model(&models.AccountEvaluation{}).Where("account_id = ?", assignID).Updates(map[string]interface{}{
			"score":        gorm.Expr("score + ?", score),
			"assign_count": gorm.Expr("assign_count + 1")}).Error != nil {
			c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
			return
		}
		c.JSON(http.StatusOK, vo.Success(nil))
	}
}

// @Tags 评价模块
// @Summary 获取账号评价
// @Description 获取账号评价
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param  accountID query uint true "账号ID"
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 500 {object} vo.ResponseVO "失败"
// @Router /evaluation/account [get]
func (ce *EvaluationController) GetAccountEvaluation(c *gin.Context) {
	accountIDStr := c.Query("accountID")
	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, vo.Fail(constants.PARAM_ERROR))
	} else {
		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, vo.Fail(constants.PARAM_ERROR))
			return
		}
		var accountEvaluationVo vo.AccountEvaluationVo
		if err = Db.Model(&models.AccountEvaluation{}).Where("account_id = ?", accountID).Take(&accountEvaluationVo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, vo.Fail(constants.ACCOUNT_NOT_EXIST))
			} else {
				c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
			}
			return
		}
		c.JSON(http.StatusOK, vo.Success(accountEvaluationVo))
	}
}
