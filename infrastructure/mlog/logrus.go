package mlog

import (
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	lfshook2 "github.com/rifflock/lfshook"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"path"
	"strconv"
	"time"
)

var logClient *logrus.Logger
var logLevels =map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"warn": logrus.WarnLevel,
	"info": logrus.InfoLevel,
}

func init(){
	logClient = logrus.New()
	logClient.Formatter = &logrus.JSONFormatter{}
}

func GetLogrus() *logrus.Logger{
	return logClient
}

func ConfigLocalFilesystemLogger(logpath,loglevel string){
	logFileName:="logs"
	baseLogPath:=path.Join(logpath,logFileName)
	level,ok:=logLevels[loglevel]
	if ok {
		logClient.SetLevel(level)
	}else{
		logClient.SetLevel(logrus.WarnLevel)
	}
	writer, err:=rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),
		rotatelogs.WithMaxAge(time.Hour*100),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err!=nil{
		logrus.Errorf("config log failed",err)
	}
	lfshook:=lfshook2.NewHook(lfshook2.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel: writer,
		logrus.WarnLevel: writer,
		logrus.PanicLevel: writer,
		logrus.ErrorLevel: writer,
	},&logrus.JSONFormatter{})
	logClient.AddHook(lfshook)
}
//每次进入网关，打印日志
func AccessBegin(c *gin.Context){
	logClient.WithFields(logrus.Fields{
		"type":"xxx",
		"host":c.Request.URL.Host,
	}).Info("xxx")
}
//每次退出网关，打印日志
func AccessEnd(c *gin.Context,time int){
	logClient.WithFields(logrus.Fields{
		"type":"xxx",
		"host":c.Request.URL.Host,
		"cost":strconv.Itoa(time),
	}).Info("xxx")
}

func GetUniId(c *gin.Context)(uid string){
	id,exits:=c.Get("requestid")
	if !exits{
		uid=uuid.NewV4().String()
	}else {
		uid=id.(string)
	}
	return uid
}