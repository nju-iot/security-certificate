package handlers

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/nju-iot/security-certificate/resp"
)

func PingV2(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "text/plain")
	_, _ = c.Writer.Write([]byte("test"))
}

func Test(c *gin.Context) *resp.JSONOutput {
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf(string(bodyBytes))
	d := EncryptData(bodyBytes)
	fmt.Printf(string(d))
	return resp.SampleJSON(c, resp.RespCodeSuccess, d)
}
