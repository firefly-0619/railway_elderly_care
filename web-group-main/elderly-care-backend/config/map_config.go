// config/map_config.go
package config

type MapConfig struct {
	AMap struct {
		Enable  bool   `yaml:"enable"`
		APIKey  string `yaml:"api_key"`
		BaseURL string `yaml:"base_url"`
	} `yaml:"amap"`
}
