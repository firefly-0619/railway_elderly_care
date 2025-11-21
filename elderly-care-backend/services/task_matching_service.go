package services

import (
	"elderly-care-backend/global"
	"elderly-care-backend/models"
	"fmt"
	"math"
)

type TaskMatchingService struct{}

// MatchVolunteersForTask 为任务匹配志愿者
func (s *TaskMatchingService) MatchVolunteersForTask(taskID uint) ([]map[string]interface{}, error) {
	var task models.Task
	if err := global.Db.First(&task, taskID).Error; err != nil {
		return nil, fmt.Errorf("获取任务失败: %v", err)
	}

	// 使用GORM查询附近志愿者
	var volunteers []map[string]interface{}
	query := `
        SELECT 
            id, nickname, latitude, longitude, address,
            (6371000 * acos(cos(radians(?)) * cos(radians(latitude)) * 
            cos(radians(longitude) - radians(?)) + sin(radians(?)) * 
            sin(radians(latitude)))) as distance
        FROM account 
        WHERE id != ? AND distance < ?
        ORDER BY distance ASC
        LIMIT 20`

	if err := global.Db.Raw(query, task.Latitude, task.Longitude, task.Latitude, task.CreatorID, 5000).Scan(&volunteers).Error; err != nil {
		return nil, err
	}

	// 计算匹配分数
	for _, volunteer := range volunteers {
		score := s.calculateMatchScore(volunteer)
		volunteer["match_score"] = score
	}

	return volunteers, nil
}

// MatchEmergencyVolunteers 紧急情况匹配
func (s *TaskMatchingService) MatchEmergencyVolunteers(taskID uint, lat, lng float64) ([]map[string]interface{}, error) {
	query := `
        SELECT 
            id, nickname, latitude, longitude, address,
            (6371000 * acos(cos(radians(?)) * cos(radians(latitude)) * 
            cos(radians(longitude) - radians(?)) + sin(radians(?)) * 
            sin(radians(latitude)))) as distance
        FROM account 
        WHERE id != ? AND distance < 3000
        ORDER BY distance ASC
        LIMIT 10`

	var volunteers []map[string]interface{}
	if err := global.Db.Raw(query, lat, lng, lat, taskID).Scan(&volunteers).Error; err != nil {
		return nil, err
	}

	return volunteers, nil
}

func (s *TaskMatchingService) calculateMatchScore(volunteer map[string]interface{}) float64 {
	distance := volunteer["distance"].(float64)
	// 距离越近分数越高
	distanceScore := math.Max(0, 100-(distance/5000)*100)
	return distanceScore
}
