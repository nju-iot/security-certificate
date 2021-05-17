package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nju-iot/security-certificate/handlers"
)

func registerRouter(r *gin.Engine) {
	r.GET("/ping", handlers.PingV2)
	// your code

}
