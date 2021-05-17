package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nju-iot/security-certificate/caller"
	dal "github.com/nju-iot/security-certificate/data"
	"github.com/nju-iot/security-certificate/logs"
	"github.com/nju-iot/security-certificate/resp"
)

type RegisterParamsV2 struct {
	Username string `form:"username" json:"username" binding:"required"`
	//Password string `form:"password" json:"password" binding:"required"`
	Email string `form:"email" json:"email"`
}

// RegisterCheckParamV2 ...
type RegisterCheckParamsV2 struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email"`
	Code     string `form:"code" json:"code" binding:"required"`
}

func RegisterV2(c *gin.Context) *resp.JSONOutput {
	// Step1. 参数校验
	params := &RegisterParamsV2{}
	err := c.Bind(&params)
	if err != nil {
		logs.Error("[Register] request-params error: params=%+v, err=%v", params, err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}

	// Step2. 查看用户/邮箱是否存在
	userInfo, dbErr := data.GetEdgexUserByNameAndEmail(params.Username, params.Email)
	// mailInfo, dbErr2 := dal.GetEdgexUserByEmail(params.Email)
	if dbErr != nil {
		logs.Error("[Register] get user failed: username=%s, err=%v", params.Username, dbErr)
		return resp.SampleJSON(c, resp.RespDatabaseError, nil)
	}
	// if dbErr2 != nil {
	// 	logs.Error("[Register] email [%s] already exists: err=%v", params.Email, dbErr2)
	// }
	if userInfo != nil {
		return resp.SampleJSON(c, resp.RespCodeUserExsit, nil)
	}

	//Step3. 发送邮箱验证码
	mailTo := []string{params.Email}
	subject := string("登录验证")
	code := randomCode()
	body := code
	err = SendMailToV2(mailTo, subject, body)
	if err != nil {
		return resp.SampleJSON(c, resp.RespCodeParamsError, "发送失败")
	}

	if checker == nil {
		checker = make(map[string]CodeCheckerV2)
	}
	var cc CodeCheckerV2
	cc.Email = params.Email
	cc.Code = code
	checker[params.Email] = cc
	return resp.SampleJSON(c, resp.RespCodeSuccess, nil)
}

// RegisterCheckV2 ...
func RegisterCheckV2(c *gin.Context) *resp.JSONOutput {
	params := &RegisterCheckParamsV2{}
	err := c.Bind(&params)
	if err != nil {
		logs.Error("[RegisterCheck] request-params error: params=%+v, err=%v", params, err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}
	//验证验证码
	err = checkCodeV2(params.Email, params.Code)
	if err != nil {
		logs.Error("[RegisterCheck] check-code error: params=%+v, err=%v", params, err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}
	//添加用户
	user := &dal.EdgexUser{
		Username:     params.Username,
		Password:     params.Password,
		Email:        params.Email,
		CreatedTime:  time.Now(),
		ModifiedTime: time.Now(),
	}
	dbErr := dal.AddEdgexUser(caller.EdgexDB, user)
	if dbErr != nil {
		return resp.SampleJSON(c, resp.RespDatabaseError, nil)
	}
	return resp.SampleJSON(c, resp.RespCodeSuccess, nil)
}
