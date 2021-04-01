package middle

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"myproject/Apigateway/config"
	"time"
)

//设置跨域请求
func Login(alias string,conf config.ConfigMap ) gin.HandlerFunc{
	return func(c *gin.Context){
		if conf[alias].JwtKey==""{
			c.Next()
			return
		}
		c.Request.Header.Set("operator","")
		auth:=c.Request.Header.Get("Authorization")
		if auth==""{
			c.Next()
			return
		}
		tmp,ok:=ParseToken(auth,conf[alias].JwtKey)
		if ok{
			claims:=tmp.(jwt.MapClaims)
			if claims!=nil{
				if operator,ok:=claims["name"];ok{
					op:=operator.(string)
					c.Request.Header.Set("operator",op)
					e:=CheckApiData(c,alias,conf[alias].AppCode,op)
					if e!=nil{
						return
					}
				}
			}
			if role,ok:=claims["role"];ok{
				c.Request.Header.Set("role",role.(string))
			}
			if interfaces, ok := claims["interfaces"]; ok {
				str, e := json.Marshal(interfaces)
				if e == nil {
					c.Request.Header.Set("interfaces", string(str))
				}
			}
			if dataRight, ok := claims["dataRight"]; ok {
				dataRight, e := json.Marshal(dataRight)
				if e == nil {
					c.Request.Header.Set("dataRight", string(dataRight))
				}
			}
		}
		c.Next()
	}
}
func ParseToken(auth,token string)(interface{},bool){
	//Parse接受令牌字符串和一个用于查找密钥的函数
	token1,err:=jwt.Parse(auth, func(token2 *jwt.Token)(interface{},error){
		if _,ok:=token2.Method.(*jwt.SigningMethodHMAC);!ok{
			return nil,fmt.Errorf("Unexpected signing method: %v", token2.Header["alg"])
		}
		return []byte(token),nil
	})
	if err!=nil{
		return "",false
	}
	if token1==nil{
		return "",false
	}
	if claims,ok:=token1.Claims.(jwt.MapClaims);ok&&int(time.Now().Unix())<int(claims["exp"].(float64)){
		return claims,true
	}
	return " ",false
}