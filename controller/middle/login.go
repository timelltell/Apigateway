package middle

import (
	"github.com/gin-gonic/gin"
	"myproject/Apigateway/config"
	"net/http"
)

//设置跨域请求
func Login(alias string,conf config.ConfigMap ) gin.HandlerFunc{
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