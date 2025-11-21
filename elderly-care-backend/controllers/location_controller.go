package controllers

import (
	"elderly-care-backend/common/constants"
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

type LocationController struct {
	amapService     *services.AMapService
	locationService *services.RealtimeLocationService
	wsService       *services.WebSocketService
	addressService  *services.AddressService
}

func NewLocationController(
	amapService *services.AMapService,
	locationService *services.RealtimeLocationService,
	wsService *services.WebSocketService,
) *LocationController {
	return &LocationController{
		amapService:     amapService,
		locationService: locationService,
		wsService:       wsService,
		addressService:  services.NewAddressService(global.Db),
	}
}

// @Tags 实时定位
// @Summary 更新用户位置
// @Description 实时更新用户当前位置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.LocationUpdate true "位置更新请求"
// @Success 200 {object} vo.ResponseVO
// @Router /location/update [post]
func (lc *LocationController) UpdateLocation(c *gin.Context) { // 首字母大写
	var update models.LocationUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusOK, vo.Fail(constants.PARAM_ERROR))
		return
	}

	// 从token中获取用户ID，确保安全
	userID := utils.GetAccountIdInContext(c)
	update.UserID = uint(userID)

	// 如果没有地址信息，使用逆地理编码获取
	if update.Address == "" && lc.amapService != nil {
		address, err := lc.amapService.ReverseGeocode(
			strconv.FormatFloat(update.Longitude, 'f', 6, 64),
			strconv.FormatFloat(update.Latitude, 'f', 6, 64),
		)
		if err == nil {
			update.Address = address.FormattedAddress
		}
	}

	if err := lc.locationService.UpdateUserLocation(update); err != nil {
		c.JSON(http.StatusOK, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	// 广播位置更新给相关用户
	if lc.wsService != nil {
		lc.wsService.BroadcastMessage("user_location_updated", map[string]interface{}{
			"user_id":   userID,
			"latitude":  update.Latitude,
			"longitude": update.Longitude,
			"address":   update.Address,
		})
	}

	c.JSON(http.StatusOK, vo.Success(gin.H{
		"user_id": userID,
		"location": gin.H{
			"lat": update.Latitude,
			"lng": update.Longitude,
		},
		"address":   update.Address,
		"accuracy":  update.Accuracy,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}))
}

// @Tags 实时定位
// @Summary 获取用户当前位置
// @Description 获取指定用户的当前位置
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param userId path int true "用户ID"
// @Success 200 {object} vo.ResponseVO{data=models.UserLocation}
// @Router /location/user/{userId} [get]
func (lc *LocationController) GetUserLocation(c *gin.Context) { // 首字母大写
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, vo.Fail(constants.PARAM_ERROR))
		return
	}

	location, err := lc.locationService.GetUserCurrentLocation(uint(userID))
	if err != nil {
		c.JSON(http.StatusOK, vo.Success(nil))
		return
	}

	c.JSON(http.StatusOK, vo.Success(location))
}

// @Tags 实时定位
// @Summary 获取附近用户
// @Description 获取附近的志愿者或求助者
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param lat query number true "纬度"
// @Param lng query number true "经度"
// @Param radius query number false "搜索半径" default(5000)
// @Param role query string false "用户角色" Enums(volunteer, user)
// @Success 200 {object} vo.ResponseVO{data=[]interface{}}
// @Router /location/nearby [get]
func (lc *LocationController) GetNearbyUsers(c *gin.Context) { // 首字母大写
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.Query("radius")
	role := c.Query("role")

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

	if role == "" {
		role = "volunteer"
	}

	users, err := lc.locationService.GetNearbyUsers(lat, lng, radius, role)
	if err != nil {
		c.JSON(http.StatusOK, vo.Fail(constants.SERVICE_ERROR))
		return
	}

	c.JSON(http.StatusOK, vo.Success(users))
}

// @Summary 计算导航路线
// @Description 根据起点终点计算步行导航路线
// @Tags 定位模块
// @Accept json
// @Produce json
// @Param from query string true "起点坐标 格式:经度,纬度"
// @Param to query string true "终点坐标 格式:经度,纬度"
// @Success 200 {object} vo.ResponseVO
// @Router /location/navigation [get]
func (lc *LocationController) CalculateNavigation(c *gin.Context) { // 首字母大写
	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, vo.Fail("起点和终点坐标不能为空"))
		return
	}

	route, err := lc.amapService.CalculateWalkingRoute(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Fail("导航计算失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.Success(route))
}

// @Summary 坐标转地址
// @Description 将经纬度坐标转换为详细地址
// @Tags 定位模块
// @Accept json
// @Produce json
// @Param lng query string true "经度"
// @Param lat query string true "纬度"
// @Success 200 {object} vo.ResponseVO
// @Router /location/reverse-geocode [get]
func (lc *LocationController) ReverseGeocode(c *gin.Context) { // 首字母大写
	lng := c.Query("lng")
	lat := c.Query("lat")

	if lng == "" || lat == "" {
		c.JSON(http.StatusBadRequest, vo.Fail("经纬度不能为空"))
		return
	}

	address, err := lc.amapService.ReverseGeocode(lng, lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Fail("地址解析失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.Success(address))
}

// @Summary 获取用户到目标的导航路线
// @Description 获取当前用户到指定目标的导航路线
// @Tags 定位模块
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param targetLat query number true "目标纬度"
// @Param targetLng query number true "目标经度"
// @Success 200 {object} vo.ResponseVO
// @Router /location/navigation/to-target [get]
func (lc *LocationController) GetNavigationToTarget(c *gin.Context) { // 首字母大写
	targetLatStr := c.Query("targetLat")
	targetLngStr := c.Query("targetLng")

	if targetLatStr == "" || targetLngStr == "" {
		c.JSON(http.StatusBadRequest, vo.Fail("目标坐标不能为空"))
		return
	}

	userID := utils.GetAccountIdInContext(c)
	currentLocation, err := lc.locationService.GetUserCurrentLocation(uint(userID))
	if err != nil || currentLocation == nil {
		c.JSON(http.StatusBadRequest, vo.Fail("请先更新您的位置信息"))
		return
	}

	from := lc.formatCoordinates(currentLocation.Longitude, currentLocation.Latitude)
	to := targetLngStr + "," + targetLatStr

	route, err := lc.amapService.CalculateWalkingRoute(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Fail("导航计算失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, vo.Success(route))
}

// @Summary 获取用户历史位置
// @Description 获取用户的最近位置记录
// @Tags 定位模块
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} vo.ResponseVO
// @Router /location/history [get]
func (lc *LocationController) GetLocationHistory(c *gin.Context) {
	userID := utils.GetAccountIdInContext(c)
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	locations, err := lc.addressService.GetUserRecentLocations(uint(userID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Fail("获取位置历史失败"))
		return
	}

	c.JSON(http.StatusOK, vo.Success(locations))
}

// @Summary 用户间导航
// @Description 计算两个用户之间的导航路线
// @Tags 定位模块
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param targetUserId query int true "目标用户ID"
// @Success 200 {object} vo.ResponseVO
// @Router /location/navigation/user [get]
func (lc *LocationController) NavigateToUser(c *gin.Context) {
	currentUserID := utils.GetAccountIdInContext(c)
	targetUserIDStr := c.Query("targetUserId")

	if targetUserIDStr == "" {
		c.JSON(http.StatusBadRequest, vo.Fail("目标用户ID不能为空"))
		return
	}

	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.Fail("目标用户ID格式错误"))
		return
	}

	currentLocation, err := lc.addressService.GetUserCurrentLocation(uint(currentUserID))
	if err != nil || currentLocation == nil {
		c.JSON(http.StatusBadRequest, vo.Fail("请先更新您的位置信息"))
		return
	}

	targetLocation, err := lc.addressService.GetUserCurrentLocation(uint(targetUserID))
	if err != nil || targetLocation == nil {
		c.JSON(http.StatusBadRequest, vo.Fail("目标用户位置不可用"))
		return
	}

	from := lc.formatCoordinates(currentLocation.Longitude, currentLocation.Latitude)
	to := lc.formatCoordinates(targetLocation.Longitude, targetLocation.Latitude)

	route, err := lc.amapService.CalculateWalkingRoute(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Fail("导航计算失败: "+err.Error()))
		return
	}

	response := map[string]interface{}{
		"navigation": route,
		"from_user": map[string]interface{}{
			"user_id":  currentUserID,
			"location": currentLocation,
			"address":  currentLocation.Address,
		},
		"to_user": map[string]interface{}{
			"user_id":  targetUserID,
			"location": targetLocation,
			"address":  targetLocation.Address,
		},
	}

	c.JSON(http.StatusOK, vo.Success(response))
}

// @Summary 根据位置ID导航
// @Description 根据位置记录ID计算导航路线
// @Tags 定位模块
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param locationId query int true "位置记录ID"
// @Success 200 {object} vo.ResponseVO
// @Router /location/navigation/location [get]
func (lc *LocationController) NavigateToLocation(c *gin.Context) {
	currentUserID := utils.GetAccountIdInContext(c)
	locationIDStr := c.Query("locationId")

	if locationIDStr == "" {
		c.JSON(http.StatusBadRequest, vo.Fail("位置ID不能为空"))
		return
	}

	locationID, err := strconv.ParseUint(locationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.Fail("位置ID格式错误"))
		return
	}

	currentLocation, err := lc.addressService.GetUserCurrentLocation(uint(currentUserID))
	if err != nil || currentLocation == nil {
		c.JSON(http.StatusBadRequest, vo.Fail("请先更新您的位置信息"))
		return
	}

	targetLocation, err := lc.addressService.GetLocationByID(uint(locationID))
	if err != nil || targetLocation == nil {
		c.JSON(http.StatusBadRequest, vo.Fail("目标位置不存在"))
		return
	}

	from := lc.formatCoordinates(currentLocation.Longitude, currentLocation.Latitude)
	to := lc.formatCoordinates(targetLocation.Longitude, targetLocation.Latitude)

	route, err := lc.amapService.CalculateWalkingRoute(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Fail("导航计算失败: "+err.Error()))
		return
	}

	response := map[string]interface{}{
		"navigation":    route,
		"from_location": currentLocation,
		"to_location":   targetLocation,
	}

	c.JSON(http.StatusOK, vo.Success(response))
}

// 辅助函数：格式化坐标
func (lc *LocationController) formatCoordinates(lng, lat float64) string {
	return strconv.FormatFloat(lng, 'f', 6, 64) + "," + strconv.FormatFloat(lat, 'f', 6, 64)
}
