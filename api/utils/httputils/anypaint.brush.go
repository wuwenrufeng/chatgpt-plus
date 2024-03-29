package httputils

import (
	"chatplus/core/types"
	"fmt"
)

type AnypaintBrush struct {
	BaseUrl    string
	RouterPath string
	AppKey     string
	AppSecret  string
}

type AnypaintData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	EnoughOk = 200
)

func (a AnypaintBrush) IsEnough(uid string) (error, *AnypaintData) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": GenToken(a.AppKey, a.AppSecret),
	}
	userUrl := a.BaseUrl + a.RouterPath + "/" + uid
	client := NewHTTPClient(3, headers) // 设置最大重试次数为3和自定义请求头
	// 发送请求
	data := &AnypaintData{}
	err := client.SendRequest(GET, userUrl, nil, data)
	if err != nil {
		logger.Error(fmt.Sprintf("刷子查询失败，%x", data), err)
		return err, nil
	}
	return nil, data
}

func (a AnypaintBrush) SubBrush(uid string, session *types.ChatSession) error {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": GenToken(a.AppKey, a.AppSecret),
	}
	userUrl := a.BaseUrl + a.RouterPath + "/" + uid
	client := NewHTTPClient(3, headers) // 设置最大重试次数为3和自定义请求头
	// 发送请求
	body := map[string]string{
		"session": session.SessionId,
		"chat_id": session.ChatId,
	}
	data := &AnypaintData{}
	err := client.SendRequest(POST, userUrl, body, data)
	return err
}
