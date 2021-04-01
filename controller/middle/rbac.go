package middle

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"sync/atomic"
)

var ApiAccessData atomic.Value
var UserRoles atomic.Value
func CheckApiData(c *gin.Context,alias,service,user_code string)error{
	data:=ApiAccessData.Load()
	apiAccessData:=data.(map[string]map[string]map[string]bool)
	url:= c.Request.URL.Path
	url=strings.TrimPrefix(url,"/"+alias+"/")
	api := strings.TrimRight(url, "/")
	//url=apiAccessData[service][api]
	//apiAccessData记录了某个服务允许哪些角色访问
	btmp:=ckeckuseraccess(user_code,apiAccessData[service][api])
	if btmp{
		return nil
	}else{
		return errors.New("无权访问")
	}
}
func ckeckuseraccess(user string,url map[string]bool)(bool){
	tmp:=UserRoles.Load()
	println(tmp)
	return true

}
