package middle

import (
	"github.com/gin-gonic/gin"
	"myproject/Apigateway/infrastructure/mlog"
	"time"
)

//中间件，实现日志打印，
func Request() gin.HandlerFunc{
	return func(c *gin.Context){
		reqid:=mlog.GetUniId(c)
		c.Set("reqeustid",reqid)
		c.Request.Header.Set("reqeustid",reqid)

		mlog.AccessBegin(c)
		beforeTime:=time.Now().UnixNano()
		c.Request.ParseForm()
		c.Next()
		endTime:=time.Now().UnixNano()
		cost:=int((endTime-beforeTime)/1000000)
		mlog.AccessEnd(c,cost)
	}
}