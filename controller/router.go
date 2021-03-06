package controller

import (
	"fmt"
	"context"
	"myproject/Apigateway/config"
	"myproject/Apigateway/controller/middle"
	"strings"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
	"myproject/Apigateway/infrastructure/mlog"
)

var engine *gin.Engine
//type Config struct {
//	Alias string `json:"alias"`
//	AppCode string `json:"app_code"`
//	JwtKey string `json:"jwt_key"`
//	Timeout int64 `json:"timeout"`
//	Url []string `json:"url"`
//	Cors bool `json:"cors"`
//}


func GetConfigFromYml() config.ConfigMap{
	var tmp config.ConfigMap
	tmp = make( config.ConfigMap )
	return tmp
}
func GetRouter() *gin.Engine{
	//初始化配置信息
	conf:=GetProxy()
	//加载全局中间件
	engine:=gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middle.Request())
	//注册路由中间件
	for alias:=range conf{
		relative:="/"+alias+"/*uri"
		engine.Any(relative,middle.SetCors(relative,conf),middle.Login(relative,conf))
	}
	return engine
}
func GetProxy() (conf1 config.ConfigMap){
	conf:=GetConfigFromYml()
	if len(conf) ==0{
		panic("no proxy config")
	}
	middle.ProxyMapDetail = make(middle.ProxyMap)
	for ali,config :=range conf{
		timeout:=config.Timeout
		urlSlice:=config.Url
		for _,singleUrl :=range urlSlice{
			urlInfo,err:=url.Parse(singleUrl)
			if err == nil{
				info:=&middle.BackendMap{
					Domain: urlInfo.Host,
					Proxy: newRVP(urlInfo,timeout),
				}
				middle.ProxyMapDetail[ali]=append(middle.ProxyMapDetail[ali],*info)
			}
		}

	}
	return

}
func newRVP(target *url.URL,timeout int64) *httputil.ReverseProxy{
	dir:=func(req *http.Request){
		req.URL.Scheme=target.Scheme
		if _,ok:=req.Header["User-Agent"];!ok{
			req.Header.Set("User-Agent","")
		}
	}
	return &httputil.ReverseProxy{
		Director: dir,
		ModifyResponse: modifyResponse,
		ErrorHandler: errHandler,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := (&net.Dialer{
					Timeout: 2000 * time.Millisecond,
					//KeepAlive: 5 * time.Second,
					//Deadline:  time.Now().Add(time.Duration(timeout) * time.Millisecond),
					DualStack: true,
				}).DialContext(ctx, network, addr)
				if err == nil {
					//超时时间设置
					conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
					conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
				}
				return conn, err
			},
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	//return &httputil.ReverseProxy{
	//	Director: dir,
	//	Transport: &http.Transport{
	//		DisableKeepAlives: true,
	//		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
	//			conn, err := (&net.Dialer{
	//				Timeout: 2000 * time.Millisecond,
	//				//KeepAlive: 5 * time.Second,
	//				//Deadline:  time.Now().Add(time.Duration(timeout) * time.Millisecond),
	//				DualStack: true,
	//			}).DialContext(ctx, network, addr)
	//			if err == nil {
	//				//超时时间设置
	//				conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	//				conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	//			}
	//			return conn, err
	//		},
	//		TLSHandshakeTimeout: 10 * time.Second,
	//	},
	//	ModifyResponse: modifyResponse,
	//	ErrorHandler:   errHandler,
	//}

}


func errHandler(cbrw http.ResponseWriter, req *http.Request, err error) {
	if strings.Index(err.Error(), "timeout") >= 0 { // backend timeout
		cbrw.WriteHeader(http.StatusGatewayTimeout)

	} else {
		cbrw.WriteHeader(http.StatusOK)
		cbrw.Write([]byte(sys.GetErrorMsg("ERR_BACKEND", err)))
	}
	//uniqueID := req.Header.Get("uniqid")
	//errStr := fmt.Sprintf("Backend Error: [%s] Request Url: [%s]", err.Error(), req.RequestURI)
	//mlog.Errorf("err:%s,requestID:%s,req:%s", errStr, uniqueID, req.URL)
}

//modify
func modifyResponse(cbrw *http.Response) error {

	//responseStatus := cbrw.StatusCode
	//switch responseStatus {
	//case 500, 502, 503, 504:
	//	cbrw.StatusCode = responseStatus + 100
	//}
	//3xx的请求
	if cbrw.StatusCode > 300 && cbrw.StatusCode < 400 {
		return nil
	}
	if cbrw.StatusCode != http.StatusOK {
		loginfo := `change http code to ` + strconv.Itoa(cbrw.StatusCode)
		uniqueID := cbrw.Header.Get("requestid")
		mlog.Infof("type:%s,requestID:%s,loginfo:%s", "reverseProxy", uniqueID, loginfo)
		return fmt.Errorf("statuscode:%d", cbrw.StatusCode)
	}
	return nil
}