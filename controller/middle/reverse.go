package middle

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"myproject/Apigateway/config"
	"net/http"
	"net/http/httputil"
)
type BackendMap struct{
	Domain string
	Proxy *httputil.ReverseProxy
}
//type ConfigMap map[string]Config
type ProxyMap map[string][]BackendMap
var ProxyMapDetail ProxyMap
func Reverse(alias string,conf config.ConfigMap) gin.HandlerFunc{
	return func(c *gin.Context){
		_,ok:=conf[alias]
		if !ok{
			c.AbortWithStatusJSON(http.StatusInternalServerError,"no config")
		}
		randBackend:=ProxyMapDetail[alias][rand.Int()%len(ProxyMapDetail[alias])]
		c.Request.URL.Path = c.Param(`uri`)
		c.Request.Host = conf[alias].Host
		c.Request.URL.Host = randBackend.Domain
		proxy:=randBackend.Proxy
		//熔断功能
		c.Next()
	}
}