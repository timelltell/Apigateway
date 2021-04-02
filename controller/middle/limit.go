package middle

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"myproject/Apigateway/config"
	"time"
)

func Limit(alias string, conf config.ConfigMap) gin.HandlerFunc {
	// create a  request/second limiter and
	// every token bucket in it will expire 1 hour after it was initially set.
	limiter := tollbooth.NewLimiter(conf[alias].Limit, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	limiter.SetMessage(`{"errno":"200429","message":"命中限速逻辑"}`)
	limiter.SetMessageContentType("application/json; charset=utf-8")

	return func(c *gin.Context) {
		if conf[alias].Limit <= 0 {
			c.Next()
			return
		}
		httpError := tollbooth.LimitByRequest(limiter, c.Writer, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, limiter.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}