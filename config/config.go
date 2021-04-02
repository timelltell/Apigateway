package config

type Config struct {
	Alias string `json:"alias"`
	AppCode string `json:"app_code"`
	JwtKey string `json:"jwt_key"`
	Timeout int64 `json:"timeout"`
	Url []string `json:"url"`
	Cors bool `json:"cors"`
	Host string `json:"host"`
	Limit float64 `json:"limit"`
}


type ConfigMap map[string]Config