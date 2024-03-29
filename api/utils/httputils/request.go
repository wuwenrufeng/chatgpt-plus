package httputils

import (
	"bytes"
	logger2 "chatplus/logger"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var logger = logger2.GetLogger()

// HTTP请求工具类
type HTTPClient struct {
	client  *http.Client
	retry   int
	headers map[string]string
}

const (
	POST   = "POST"
	PUT    = "PUT"
	GET    = "GET"
	DELETE = "DELETE"
)

// 创建HTTP请求工具类
func NewHTTPClient(retry int, headers map[string]string) *HTTPClient {
	return &HTTPClient{
		client:  &http.Client{},
		retry:   retry,
		headers: headers,
	}
}

// 发送HTTP请求
func (c *HTTPClient) SendRequest(method, url string, body interface{}, data interface{}) error {
	var resp *http.Response
	var err error
	for i := 0; i <= c.retry; i++ {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return fmt.Errorf("创建请求失败: %v", err)
		}

		// 设置请求头
		for key, value := range c.headers {
			req.Header.Set(key, value)
		}

		// 设置请求体
		if body != nil {
			marshal, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("body err: %v", err)
			}
			req.Body = ioutil.NopCloser(bytes.NewBuffer(marshal))
		}

		resp, err = c.client.Do(req)
		if err == nil {
			break
		}

		if i < c.retry {
			logger.Info(fmt.Sprintf("发送请求失败，%s,进行第 %d 次重试...\n", url, i+1))
			time.Sleep(time.Second) // 可以根据需要调整重试间隔
		}
	}
	if resp == nil {
		return fmt.Errorf("发送请求失败，已达到最大重试次数")
	}
	logger.Infof("http response code:%d ,url:(%s) \n", resp.StatusCode, url)
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}
	return json.Unmarshal(respByte, data)
}

func GenToken(appKey string, appSecret string) string {
	currentTime := strconv.FormatInt(time.Now().Unix(), 10)
	data := appKey + appSecret + currentTime

	hash := sha256.Sum256([]byte(data))
	hexDigest := fmt.Sprintf("%x", hash)

	b64str := fmt.Sprintf("%s:%s:%s", appKey, currentTime, hexDigest)
	payload := base64.StdEncoding.EncodeToString([]byte(b64str))

	token := fmt.Sprintf("Secret %s", payload)

	return token
}
