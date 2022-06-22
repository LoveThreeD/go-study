package response

import (
	"encoding/json"
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	prefix     = "gin-admin"
	ReqBodyKey = prefix + "/req-body"
	ResBodyKey = prefix + "/res-body"
)

type Response struct {
	Code int         `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 错误描述
	Data interface{} `json:"data"` // 返回数据
}

// WithMsg 自定义响应信息
func (res *Response) WithMsg(message string) Response {
	return Response{
		Code: res.Code,
		Msg:  message,
		Data: res.Data,
	}
}

// WithData 追加响应数据
func (res *Response) WithData(data interface{}) Response {
	return Response{
		Code: res.Code,
		Msg:  res.Msg,
		Data: data,
	}
}

// ToString 返回 JSON 格式的错误详情
func (res *Response) ToString() string {
	err := &struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{
		Code: res.Code,
		Msg:  res.Msg,
		Data: res.Data,
	}
	raw, _ := json.Marshal(err)
	return string(raw)
}

// ResJSON Response json data with status code
func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		logger.Fatal(err)
	}
	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

// ResSuccess Response data object
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, OK.WithData(v))
}

func ResSuccessWithData(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

func ResFailedWithData(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

func ResFailed(c *gin.Context) {
	ResJSON(c, http.StatusOK, Err)
}

// 构造函数
func response(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
