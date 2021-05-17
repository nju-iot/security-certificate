package main

import (
	"flag"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/nju-iot/security-certificate/caller"
	"github.com/nju-iot/security-certificate/config"
	"github.com/nju-iot/security-certificate/logs"

	"github.com/nju-iot/security-certificate/middleware/cors"
	"github.com/nju-iot/security-certificate/middleware/session"
	"go.uber.org/zap"
)

func main() {

	// r := gin.New()
	var confFilePath string
	flag.StringVar(&confFilePath, "conf", "", "Specify local configuration file path")
	flag.Parse()

	config.LoadConfig(confFilePath)
	logs.InitLogs()
	caller.InitClient()

	gin.SetMode(config.Server.RunMode)

	r := gin.New()
	r.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(zap.L(), true))

	// 允许跨域访问
	r.Use(cors.MiddlewareCors())

	// session中间件
	r.Use(session.EnableRedisSession())
	r.Use(session.MiddlewareSession())

	registerRouter(r)
	_ = r.Run(config.Server.Port)
}
