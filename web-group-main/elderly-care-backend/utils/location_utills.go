// utils/location_utils.go
package utils

import (
	"fmt"
	"math"
)

// CalculateDistance 计算两点间距离（米）- Haversine公式
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000 // 地球半径(米)

	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	Δφ := (lat2 - lat1) * math.Pi / 180
	Δλ := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// CalculateBearing 计算方位角（度）
func CalculateBearing(lat1, lng1, lat2, lng2 float64) float64 {
	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	λ1 := lng1 * math.Pi / 180
	λ2 := lng2 * math.Pi / 180

	y := math.Sin(λ2-λ1) * math.Cos(φ2)
	x := math.Cos(φ1)*math.Sin(φ2) - math.Sin(φ1)*math.Cos(φ2)*math.Cos(λ2-λ1)
	θ := math.Atan2(y, x)

	bearing := (θ*180/math.Pi + 360)
	return math.Mod(bearing, 360)
}

// FormatDistance 格式化距离显示
func FormatDistance(distance float64) string {
	if distance < 1000 {
		return fmt.Sprintf("%.0f米", distance)
	}
	return fmt.Sprintf("%.1f公里", distance/1000)
}

// FormatDuration 格式化时间显示
func FormatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%d秒", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%d分钟", seconds/60)
	}
	return fmt.Sprintf("%d小时%d分钟", seconds/3600, (seconds%3600)/60)
}
