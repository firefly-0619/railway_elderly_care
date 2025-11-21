// services/amap_service.go
package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type AMapService struct {
	APIKey  string
	BaseURL string
}

// 高德API响应结构
type AMapRouteResponse struct {
	Status string `json:"status"`
	Info   string `json:"info"`
	Count  string `json:"count"`
	Route  Route  `json:"route"`
}

type Route struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Paths       []Path `json:"paths"`
}

// 修改这些字段从 int 为 string
type Path struct {
	Distance     string      `json:"distance"` // 改为 string
	Duration     string      `json:"duration"` // 改为 string
	Strategy     string      `json:"strategy"`
	Tolls        float64     `json:"tolls"`
	TollDistance int         `json:"toll_distance"`
	Steps        []AMapStep  `json:"steps"`
	Restriction  interface{} `json:"restriction"`
}

type AMapStep struct {
	Instruction interface{} `json:"instruction"`
	Orientation interface{} `json:"orientation"`
	Road        interface{} `json:"road"`
	Distance    interface{} `json:"distance"`
	Tolls       interface{} `json:"tolls"`
	Polyline    string      `json:"polyline"`
	Action      interface{} `json:"action"`
}

// 我们的通用响应结构
type RouteResponse struct {
	Distance    int    `json:"distance"`    // 米
	Duration    int    `json:"duration"`    // 秒
	Polyline    string `json:"polyline"`    // 路线坐标
	Steps       []Step `json:"steps"`       // 导航步骤
	Origin      string `json:"origin"`      // 起点坐标
	Destination string `json:"destination"` // 终点坐标
}

type Step struct {
	Instruction string `json:"instruction"` // 导航指令
	Distance    int    `json:"distance"`    // 步骤距离
	Duration    int    `json:"duration"`    // 步骤时间
	Road        string `json:"road"`        // 道路名称
	Action      string `json:"action"`      // 动作类型
	Polyline    string `json:"polyline"`    // 步骤路线
}

// 逆地理编码响应
type AMapGeocodeResponse struct {
	Status    string    `json:"status"`
	Info      string    `json:"info"`
	Regeocode Regeocode `json:"regeocode"`
}

type Regeocode struct {
	AddressComponent AddressComponent `json:"addressComponent"`
	FormattedAddress string           `json:"formatted_address"`
}

type AddressComponent struct {
	Province     string       `json:"province"`
	City         string       `json:"city"`
	District     string       `json:"district"`
	Township     string       `json:"township"`
	Neighborhood Neighborhood `json:"neighborhood"`
	Building     Building     `json:"building"`
	Street       Street       `json:"street"`
}

type Neighborhood struct {
	Name string `json:"name"`
}

type Building struct {
	Name string `json:"name"`
}

type Street struct {
	Name string `json:"name"`
}

type AddressInfo struct {
	FormattedAddress string `json:"formatted_address"`
	Province         string `json:"province"`
	City             string `json:"city"`
	District         string `json:"district"`
	Township         string `json:"township"`
	Street           string `json:"street"`
	Building         string `json:"building"`
	Neighborhood     string `json:"neighborhood"`
}

// 计算步行路线
func (a *AMapService) CalculateWalkingRoute(origin, destination string) (*RouteResponse, error) {
	baseURL := a.BaseURL + "/v3/direction/walking"
	params := url.Values{}
	params.Add("key", a.APIKey)
	params.Add("origin", origin) // 格式: "经度,纬度"
	params.Add("destination", destination)
	params.Add("output", "JSON")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("请求高德API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var amapResp AMapRouteResponse
	if err := json.Unmarshal(body, &amapResp); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	if amapResp.Status != "1" {
		return nil, fmt.Errorf("高德API错误: %s", amapResp.Info)
	}

	if len(amapResp.Route.Paths) == 0 {
		return nil, fmt.Errorf("未找到可行路线")
	}

	path := amapResp.Route.Paths[0]
	return a.convertToRouteResponse(path, origin, destination), nil
}

// 坐标转地址（逆地理编码）
func (a *AMapService) ReverseGeocode(lng, lat string) (*AddressInfo, error) {
	baseURL := a.BaseURL + "/v3/geocode/regeo"
	params := url.Values{}
	params.Add("key", a.APIKey)
	params.Add("location", lng+","+lat)
	params.Add("output", "JSON")
	params.Add("extensions", "base")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("请求高德逆地理编码失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	// 添加调试输出
	fmt.Printf("高德API原始响应: %s\n", string(body))

	var geocodeResp AMapGeocodeResponse
	if err := json.Unmarshal(body, &geocodeResp); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	if geocodeResp.Status != "1" {
		return nil, fmt.Errorf("高德逆地理编码错误: %s", geocodeResp.Info)
	}

	return &AddressInfo{
		FormattedAddress: geocodeResp.Regeocode.FormattedAddress,
		Province:         geocodeResp.Regeocode.AddressComponent.Province,
		City:             geocodeResp.Regeocode.AddressComponent.City,
		District:         geocodeResp.Regeocode.AddressComponent.District,
		Township:         geocodeResp.Regeocode.AddressComponent.Township,
		Street:           geocodeResp.Regeocode.AddressComponent.Street.Name,
		Building:         geocodeResp.Regeocode.AddressComponent.Building.Name,
		Neighborhood:     geocodeResp.Regeocode.AddressComponent.Neighborhood.Name,
	}, nil
}

// 转换高德响应到通用格式
// 修改转换函数来处理所有可能的类型
func (a *AMapService) convertToRouteResponse(path Path, origin, destination string) *RouteResponse {
	steps := make([]Step, len(path.Steps))

	// 安全转换距离和时长
	distance := a.safeStringToInt(path.Distance)
	duration := a.safeStringToInt(path.Duration)

	for i, amapStep := range path.Steps {
		stepDistance := a.safeStringToInt(amapStep.Distance)

		steps[i] = Step{
			Instruction: a.toString(amapStep.Instruction),
			Distance:    stepDistance,
			Duration:    stepDistance / 80, // 估算时间
			Road:        a.toString(amapStep.Road),
			Action:      a.toString(amapStep.Action),
			Polyline:    amapStep.Polyline,
		}
	}

	return &RouteResponse{
		Distance:    distance,
		Duration:    duration,
		Polyline:    a.mergePolylines(path.Steps),
		Steps:       steps,
		Origin:      origin,
		Destination: destination,
	}
}

// 辅助函数：安全转换为字符串
func (a *AMapService) toString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if str, ok := v[0].(string); ok {
				return str
			}
		}
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// 辅助函数：安全转换为整数
func (a *AMapService) safeStringToInt(value interface{}) int {
	str := a.toString(value)
	if str == "" {
		return 0
	}
	result, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return result
}

// 合并所有步骤的路线
func (a *AMapService) mergePolylines(steps []AMapStep) string {
	var polylines []string
	for _, step := range steps {
		polylines = append(polylines, step.Polyline)
	}
	// 简单合并，实际可能需要更复杂的处理
	if len(polylines) > 0 {
		return polylines[0]
	}
	return ""
}

// 计算两点间直线距离（米）
func (a *AMapService) CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// 使用Haversine公式计算大圆距离
	const R = 6371000 // 地球半径(米)

	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	Δφ := (lat2 - lat1) * math.Pi / 180
	Δλ := (lng2 - lng1) * math.Pi / 180

	aVal := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(aVal), math.Sqrt(1-aVal))

	return R * c
}
