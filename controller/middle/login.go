package middle

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"myproject/Apigateway/config"
	"net/http"
)

//设置跨域请求
func Login(alias string,conf config.ConfigMap ) gin.HandlerFunc{
	return func(c *gin.Context){
		if conf[alias].JwtKey==""{
			c.Next()
			return
		}
		auth
	}
}
func ParseToken(auth,token string)(interface{},bool){
	token1,err:=jwt.Parse(auth, func(token2 *jwt.Token)(interface{},error){
		if _,ok:=token2.Method.(*jwt.SigningMethodHMAC);!ok{
			return nil,fmt.Errorf("failed")
		}
		return []byte(token),nil
	})
	if err!=nil{
		return "",err
	}
	return " ",nil
}