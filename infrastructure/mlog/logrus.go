package mlog

import "github.com/sirupsen/logrus"

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