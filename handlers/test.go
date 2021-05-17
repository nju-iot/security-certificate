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
	fmt.Printf("原始数据: %s\n", string(bodyBytes))
	d := testEncryptData(bodyBytes)
	fmt.Printf("公钥加密后的数据：%s\n", d)

	fmt.Printf("%s\n", string(DecryptData(d)))
	return resp.SampleJSON(c, resp.RespCodeSuccess, string(DecryptData(d)))
}
