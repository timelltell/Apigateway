package middle

import (
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"math/rand"
	"myproject/Apigateway/config"
	"net/http"
	"net/http/httputil"
	"strconv"
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
		output:=make(chan bool,1)
		var errChan chan error
		errChan=hystrix.Go(alias,func()error{
			defer func(){
				if r:=recover();r!=nil{
					errChan<-errors.New("error panic")
					return
				}
			}()
			proxy.ServeHTTP(c.Writer,c.Request)
			if c.Writer.Status()>=500{
				er:=errors.New(strconv.Itoa(c.Writer.Status()))
				errChan<-er
				return er
			}
			output<-true
			return nil
		},nil)
		select {
		case <-output:
			//ok
		case <-errChan:
			c.AbortWithStatusJSON(http.StatusInternalServerError,"熔断")
		}
		c.Next()
	}
}