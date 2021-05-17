package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nju-iot/security-certificate/handlers"
	"github.com/nju-iot/security-certificate/resp"
)

func registerRouter(r *gin.Engine) {
	r.GET("/ping", handlers.PingV2)
	// your code
	userRouter := r.Group("security-certificate/user")
	{
		userRouter.POST("/test", resp.JSONOutPutWrapper(handlers.Test))
	}

}
