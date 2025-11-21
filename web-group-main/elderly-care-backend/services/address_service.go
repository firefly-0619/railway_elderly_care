// services/address_service.go
package services

import (
	"elderly-care-backend/models"

	"gorm.io/gorm"
)

type AddressService struct {
	db *gorm.DB
}

func NewAddressService(db *gorm.DB) *AddressService {
	return &AddressService{db: db}
}

// GetUserRecentLocations 获取用户最近的定位记录
func (s *AddressService) GetUserRecentLocations(userID uint, limit int) ([]models.UserLocation, error) {
	var locations []models.UserLocation
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&locations).Error
	return locations, err
}

// GetUserCurrentLocation 获取用户当前位置
func (s *AddressService) GetUserCurrentLocation(userID uint) (*models.UserLocation, error) {
	var location models.UserLocation
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&location).Error

	if err != nil {
		return nil, err
	}
	return &location, nil
}

// GetLocationByID 根据ID获取特定位置记录
func (s *AddressService) GetLocationByID(locationID uint) (*models.UserLocation, error) {
	var location models.UserLocation
	err := s.db.Where("id = ?", locationID).First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}
