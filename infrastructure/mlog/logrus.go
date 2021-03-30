package mlog

import (
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	lfshook2 "github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path"
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

func AccessBegin(c *gin.Context){
	logClient.WithFields(logrus.Fields{
		"type":"xxx",
		"host":c.Request.URL.Host,
	}).Info("xxx")
}