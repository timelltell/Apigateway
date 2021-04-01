package middle

import (
	"github.com/gin-gonic/gin"
	"myproject/Apigateway/config"
	"net/http"
)
//type Config struct {
//	Alias string `json:"alias"`
//	AppCode string `json:"app_code"`
//	JwtKey string `json:"jwt_key"`
//	Timeout int64 `json:"timeout"`
//	Url []string `json:"url"`
//	Cors bool `json:"cors"`
//}
//
//type ConfigMap map[string]Config
//设置跨域请求
func SetCors(alias string,conf config.ConfigMap ) gin.HandlerFunc{
	return func(c *gin.Context){
		if conf[alias].Cors{
			meth:=c.Request.Method
			c.Header("Access-Control-Allow-Methods", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Expose-Headers", "*")
			c.Header("Access-Control-Allow-Headers", "*")
			if meth=="OPTIONS"{
				c.AbortWithStatus(http.StatusNoContent)
			}
		}
		c.Next()
	}
}