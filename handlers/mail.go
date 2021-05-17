package handlers

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nju-iot/security-certificate/logs"
	"github.com/nju-iot/security-certificate/resp"

	"gopkg.in/gomail.v2"
)

type MailParam struct {
	Email string `form:"email" json:"email" binding:"required"`
}

func SendMailToV2(mailTo []string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码

	//mailConn := map[string]string{
	//  "user": "xxx@163.com",
	//  "pass": "your password",
	//  "host": "smtp.163.com",
	//  "port": "465",
	//}

	mailConn := map[string]string{
		"user": "2369351080@qq.com",
		"pass": "inkdesahnqrjdjeg",
		"host": "smtp.qq.com",
		"port": "465",
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mailConn["user"], "NJU-IOT-EDGEX验证")) //这种方式可以添加别名，即“XX官方”　 //说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}

func SendMailV2(c *gin.Context) *resp.JSONOutput {
	params := &MailParam{}
	err := c.Bind(&params)
	if err != nil {
		logs.Error("[SendMail] request-params error: params=%+v, err=%v", params, err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}
	mailTo := []string{params.Email}
	subject := string("登录验证")
	code := randomCode()
	body := code
	err = SendMailToV2(mailTo, subject, body)
	if err != nil {
		return resp.SampleJSON(c, resp.RespCodeParamsError, "发送失败")
	}
	return resp.SampleJSON(c, resp.RespCodeSuccess, nil)
}

func randomCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return vcode
}

type CodeCheckerV2 struct {
	Email string
	Code  string
}

var checker map[string]CodeCheckerV2

func checkCodeV2(email string, code string) error {
	if checker == nil {
		checker = make(map[string]CodeCheckerV2)
		return errors.New("no email found")
	}
	c, ok := checker[email]
	if !ok {
		return errors.New("no email found")
	}
	if c.Code == code {
		return nil
	}
	return errors.New("wrong code")
}
