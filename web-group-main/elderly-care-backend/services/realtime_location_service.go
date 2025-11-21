// services/realtime_location_service.go
package services

import (
	"elderly-care-backend/models"
	"elderly-care-backend/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type RealtimeLocationService struct {
	db *gorm.DB
}

// 修改构造函数，接收db依赖
func NewRealtimeLocationService(db *gorm.DB) *RealtimeLocationService {
	return &RealtimeLocationService{db: db}
}

// UpdateUserLocation 更新用户实时位置
func (s *RealtimeLocationService) UpdateUserLocation(update models.LocationUpdate) error {
	// 保存到数据库
	location := models.UserLocation{
		UserID:    update.UserID,
		Latitude:  update.Latitude,
		Longitude: update.Longitude,
		Address:   update.Address,
		Accuracy:  update.Accuracy,
		Speed:     update.Speed,
		Heading:   update.Heading,
	}

	if err := s.db.Create(&location).Error; err != nil {
		return err
	}

	// 同时更新用户表的最后位置（用于快速查询）
	if err := s.db.Model(&models.Account{}).
		Where("id = ?", update.UserID).
		Updates(map[string]interface{}{
			"latitude":             update.Latitude,
			"longitude":            update.Longitude,
			"address":              update.Address,
			"last_location_update": time.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}

// GetUserCurrentLocation 获取用户当前位置
func (s *RealtimeLocationService) GetUserCurrentLocation(userID uint) (*models.UserLocation, error) {
	var location models.UserLocation
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&location).Error

	if err != nil {
		return nil, err
	}
	return &location, nil
}

// GetNearbyUsers 获取附近用户（用于匹配和显示）
func (s *RealtimeLocationService) GetNearbyUsers(lat, lng, radius float64, userType string) ([]map[string]interface{}, error) {
	fmt.Printf("=== 附近用户查询详细调试 ===\n")
	fmt.Printf("查询角色: '%s'\n", userType)

	// 先单独检查account表的角色数据
	var accounts []models.Account
	s.db.Where("role = ?", userType).Find(&accounts)
	fmt.Printf("account表中角色为'%s'的用户数量: %d\n", userType, len(accounts))
	for _, acc := range accounts {
		fmt.Printf(" - 用户%d: %s (角色: %s)\n", acc.ID, acc.Nickname, acc.Role)
	}

	var userLocations []models.UserLocation
	timeLimit := time.Now().Add(-24 * time.Hour)

	// 执行原始查询
	err := s.db.Joins("JOIN account a ON user_location.user_id = a.id").
		Where("user_location.created_at > ?", timeLimit).
		Where("a.role = ?", userType).
		Order("user_location.created_at DESC").
		Find(&userLocations).Error

	fmt.Printf("JOIN查询结果: %d 条记录\n", len(userLocations))

	// 如果没有结果，尝试不加角色过滤
	if len(userLocations) == 0 {
		fmt.Printf("=== 尝试不加角色过滤 ===\n")
		var allLocations []models.UserLocation
		s.db.Joins("JOIN account a ON user_location.user_id = a.id").
			Where("user_location.created_at > ?", timeLimit).
			Order("user_location.created_at DESC").
			Find(&allLocations)
		fmt.Printf("不加角色过滤的结果: %d 条记录\n", len(allLocations))

		for _, loc := range allLocations {
			var acc models.Account
			s.db.First(&acc, loc.UserID)
			fmt.Printf(" - 用户%d: 角色=%s, 位置=%f,%f\n", loc.UserID, acc.Role, loc.Latitude, loc.Longitude)
		}
	}

	fmt.Printf("=== 附近用户查询调试 ===\n")
	fmt.Printf("中心点: %f, %f\n", lat, lng)
	fmt.Printf("半径: %f米, 角色: %s\n", radius, userType)

	if err != nil {
		fmt.Printf("查询错误: %v\n", err)
		return nil, err
	}

	fmt.Printf("找到%d个符合条件的用户位置记录\n", len(userLocations))

	// 显示每个找到的用户
	for i, loc := range userLocations {
		fmt.Printf("用户%d: ID=%d, 位置=%f,%f, 时间=%s\n",
			i+1, loc.UserID, loc.Latitude, loc.Longitude, loc.CreatedAt.Format("15:04:05"))
	}

	// 计算距离并过滤
	var nearbyUsers []map[string]interface{}
	for _, location := range userLocations {
		distance := utils.CalculateDistance(lat, lng, location.Latitude, location.Longitude)
		fmt.Printf("用户%d距离: %.2f米\n", location.UserID, distance)

		if distance <= radius {
			user := map[string]interface{}{
				"id":          location.UserID,
				"latitude":    location.Latitude,
				"longitude":   location.Longitude,
				"address":     location.Address,
				"distance":    distance,
				"last_update": location.CreatedAt,
			}
			nearbyUsers = append(nearbyUsers, user)
			fmt.Printf("用户%d在范围内，添加到结果\n", location.UserID)
		} else {
			fmt.Printf("用户%d超出范围\n", location.UserID)
		}
	}

	fmt.Printf("最终找到%d个附近用户\n", len(nearbyUsers))
	return nearbyUsers, nil
}

// CalculateRealTimeNavigation 实时导航计算
func (s *RealtimeLocationService) CalculateRealTimeNavigation(startLat, startLng, endLat, endLng float64) (map[string]interface{}, error) {
	// 计算实时距离和方向
	distance := utils.CalculateDistance(startLat, startLng, endLat, endLng)

	// 计算方位角
	bearing := utils.CalculateBearing(startLat, startLng, endLat, endLng)

	// 估算到达时间（假设步行速度1.4m/s）
	eta := int(distance / 1.4)

	return map[string]interface{}{
		"distance":    distance,
		"bearing":     bearing,
		"eta_seconds": eta,
		"start_point": map[string]float64{"lat": startLat, "lng": startLng},
		"end_point":   map[string]float64{"lat": endLat, "lng": endLng},
		"update_time": time.Now(),
	}, nil
}
