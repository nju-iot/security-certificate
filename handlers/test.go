package handlers

import (
	"github.com/gin-gonic/gin"
)

func PingV2(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "text/plain")
	_, _ = c.Writer.Write([]byte("test"))
}
