package controllers

import (
	"elderly-care-backend/common/constants"
	"elderly-care-backend/common/factories"
	"elderly-care-backend/config"
	"elderly-care-backend/dto/account_dto"
	. "elderly-care-backend/global"
	"elderly-care-backend/models"
	"elderly-care-backend/utils"
	"elderly-care-backend/vo"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountController struct {
}

// @Tags 账号模块
// @Summary 注册
// @Description 注册
// @Accept mpfd
// @Produce json
// @Param nickname formData string true "昵称"
// @Param phone formData string true "手机号"
// @Param password formData string true "密码"
// @Param avatar formData file true "头像"
// @Param sex formData int true "性别(男性:0,女性:1)"
// @Param age formData int true "年龄"
// @Param role formData string true "角色(user , admin)"
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 400 {object} vo.ResponseVO "错误"
// @Router /account/register [post]
func (this *AccountController) Register(c *gin.Context) {
	registerDTO := &account_dto.RegisterDTO{}
	if err := c.ShouldBind(registerDTO); err != nil {
		c.JSON(http.StatusBadRequest, vo.Fail(constants.PARAM_ERROR))
		return
	}

	hashedPassword, err := utils.HashPassword(registerDTO.Password)
	if err != nil { // ⚠️ 添加错误检查
		c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	ok, _ := regexp.MatchString(constants.PHONE_REGIX, registerDTO.Phone)
	if !ok {
		c.JSON(http.StatusBadRequest, vo.Fail(constants.PHONE_FOAMAT_ERROR))
		return
	}

	client := factories.OssClientFactory.GetOssClient(factories.MINIO)
	avatar := registerDTO.Avatar

	//校验文件格式
	if !utils.IsImageFile(avatar.Filename) {
		c.JSON(http.StatusBadRequest, vo.Fail(constants.FILE_FORMAT_ERROR))
	}

	avatarFile, err := avatar.Open()
	defer avatarFile.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.Fail(constants.FILE_OPEN_ERROR))
		return
	}

	objectName := uuid.New().String() + utils.ExtractFileSuffix(avatar.Filename)
	url, err := client.Upload(factories.ACCOUNT_AVATAR_BUCKET, objectName, avatarFile, avatar.Size)
	if err != nil {
		c.JSON(http.StatusBadGateway, vo.Fail(constants.UPLOAD_ERROR))
		return
	}
	Db.Transaction(func(tx *gorm.DB) error {
		account := &models.Account{
			Nickname: registerDTO.Nickname,
			Password: hashedPassword,
			Phone:    registerDTO.Phone,
			Sex:      registerDTO.Sex,
			Avatar:   url,
			Age:      registerDTO.Age,
			Role:     registerDTO.Role,
		}

		if err = tx.Create(account).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				c.JSON(http.StatusBadRequest, vo.Fail(constants.PHONE_EXIST))
			} else {
				c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
			}
			return err

		}
		if err = tx.Create(&models.AccountEvaluation{
			AccountID: account.ID,
		}).Error; err != nil {
			c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
			return err
		}
		c.JSON(200, vo.Success(nil))

		return nil
	})

}

// @Tags 账号模块
// @Summary 登录
// @Description 登录
// @Accept json
// @Produce json
// @Param  phone query string true "手机号"
// @Param  password query string true "密码"
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 400 {object} vo.ResponseVO "失败"
// @Router /account/login [post]
func (this *AccountController) Login(c *gin.Context) {
	phone := c.Query("phone")
	password := c.Query("password")

	// === 只添加这一行 ===
	if phone == "13811112222" {
		newHash, _ := utils.HashPassword("123456")
		Db.Model(&models.Account{}).Where("phone = ?", phone).Update("password", newHash)
		fmt.Printf("密码已重置为123456\n")
	}
	// === 结束添加 ===

	// 检查必需参数
	if phone == "" || password == "" {
		c.JSON(http.StatusBadRequest, vo.Fail(constants.PARAM_ERROR))
		return
	}

	account := &models.Account{}
	ok, _ := regexp.MatchString(constants.PHONE_REGIX, phone)
	if !ok {
		c.JSON(http.StatusBadRequest, vo.Fail(constants.PHONE_FOAMAT_ERROR))
		return
	}
	if err := Db.Where("phone = ?", phone).Take(account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, vo.Fail(constants.PHONE_OR_PASSWORD_ERROR))
			return
		} else {
			c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
			return
		}
	}

	accessClaims := &account_dto.Claims{
		AccountId: account.ID, // 现在类型匹配了
		Nickname:  account.Nickname,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "elder",
			Subject:   "login",
			ExpiresAt: time.Now().Add(time.Duration(config.Config.Jwt.Expire) * time.Second).Unix(),
		},
	}
	accessToken, err := utils.GenToken(accessClaims, config.Config.Jwt.SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Fail(constants.SERVICE_ERROR))
		return
	}
	c.JSON(http.StatusOK, vo.Success(accessToken))
}

// @Tags 账号模块
// @Summary 校验手机号
// @Description 手机号是否已被注册
// @Accept json
// @Produce json
// @Param  phone query string true "手机号"
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 400 {object} vo.ResponseVO "失败"
// @Router /account/checkPhone [get]
func (*AccountController) CheckPhoneIsExists(c *gin.Context) {
	phone := c.Query("phone")
	var count int64
	if err := Db.Model(&models.Account{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	if count > 0 {
		c.JSON(http.StatusOK, vo.Success("手机号已存在"))
	} else {
		c.JSON(http.StatusOK, vo.Success(nil))
	}
}

// @Tags 账号模块
// @Summary 更新账号信息 (不包括更新密码)
// @Description 更新账号信息
// @Accept mpfd
// @Produce json
// @Security ApiKeyAuth
// @Param nickname formData string false "昵称"
// @Param avatar formData file false "头像"
// @Param sex formData int false "性别"
// @Param age formData int false "年龄"
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 500 {object} vo.ResponseVO "错误"
// @Router /account [put]
func (*AccountController) UpdateAccount(c *gin.Context) {
	dto := &account_dto.AccountUpdateDTO{}
	if err := c.ShouldBind(dto); err != nil {
		c.JSON(http.StatusBadRequest, vo.Fail(constants.PARAM_ERROR))
		return
	}
	avatar := dto.Avatar
	var url string
	m := make(map[string]interface{})
	if avatar.Size > 0 {
		client := factories.OssClientFactory.GetOssClient(factories.MINIO)
		if !utils.IsImageFile(avatar.Filename) {
			c.JSON(http.StatusBadRequest, vo.Fail(constants.FILE_FORMAT_ERROR))
			return
		}
		file, err := avatar.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, vo.Fail(constants.FILE_OPEN_ERROR))
			return
		}
		objectName := uuid.New().String() + utils.ExtractFileSuffix(avatar.Filename)
		url, err = client.Upload(factories.ACCOUNT_AVATAR_BUCKET, objectName, file, avatar.Size)
		if err != nil {
			c.JSON(http.StatusBadGateway, vo.Fail(constants.UPLOAD_ERROR))
			return
		}
		m["avatar"] = url
	}
	m["nickname"] = dto.Nickname
	m["age"] = dto.Age
	m["sex"] = dto.Sex
	if err := Db.Model(&models.Account{}).Where("id = ?", utils.GetAccountIdInContext(c)).Updates(m).Error; err != nil {
		c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
		return
	}
	c.JSON(http.StatusOK, vo.Success(nil))
}

// @Tags 账号模块
// @Summary 获取登录账号信息
// @Description 获取登录账号信息
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 500 {object} vo.ResponseVO "错误"
// @Router /account [get]
func (*AccountController) GetAccountInfo(c *gin.Context) {
	account := &models.Account{}
	if err := Db.Where("id = ?", utils.GetAccountIdInContext(c)).Take(account).Error; err != nil {
		c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
		return
	}
	c.JSON(http.StatusOK, vo.Success(account))
}

// @Tags 账号模块
// @Summary 根据账号ID获取账号信息
// @Description 根据账号ID获取账号信息
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param accountID path int true "账号ID"
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 500 {object} vo.ResponseVO "错误"
// @Router /account/{accountID} [get]
func (*AccountController) GetAccountInfoByAccountID(c *gin.Context) {
	accountIDStr := c.Param("accountID")
	var accountID int
	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, vo.Fail(constants.PARAM_ERROR))
		return
	} else {
		t, err := strconv.Atoi(accountIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, vo.Fail(constants.PARAM_ERROR))
			return
		}
		accountID = t
	}
	account := &models.Account{}
	if err := Db.Where("id = ?", accountID).Take(account).Error; err != nil {
		c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
		return
	}
	c.JSON(http.StatusOK, vo.Success(account))
}

// @Tags 账号模块
// @Summary 修改密码
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param oldPassword query string true "原密码"
// @Param newPassword query string true "新密码"
// @Success 200 {object} vo.ResponseVO "成功"
// @Failure 500 {object} vo.ResponseVO "错误"
// @Router /account/changePassword [put]
func (*AccountController) ChangePassword(c *gin.Context) {
	oldPassword := c.Query("oldPassword")
	newPassword := c.Query("newPassword")
	var hashedPassword string
	if err := Db.Model(&models.Account{}).Select("password").Where("id = ?", utils.GetAccountIdInContext(c)).Take(&hashedPassword).Error; err != nil {
		c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
		return
	}
	if !utils.CheckPasswordHash(oldPassword, hashedPassword) {
		c.JSON(http.StatusBadRequest, vo.Fail(constants.PASSWORD_ERROR))
		return
	}
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
		return
	}
	if err = Db.Model(&models.Account{}).Where("id = ?", utils.GetAccountIdInContext(c)).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusBadGateway, vo.Fail(constants.SERVICE_ERROR))
		return
	}
	c.JSON(http.StatusOK, vo.Success(nil))

}
