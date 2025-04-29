package btapi

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/carlmjohnson/requests"
)

type BtApiClient struct {
	btUrl              string
	btKey              string
	ignoreSslTransport *http.Transport
}

// 生成 request_token: md5(request_time + md5(api_sk))
func (bt *BtApiClient) GenerateRequestToken(requestTime time.Time) string {
	timestamp := fmt.Sprintf("%d", requestTime.Unix()) // Unix 时间戳（秒）
	innerMd5 := md5.Sum([]byte(bt.btKey))              // md5(api_sk)
	innerMd5Str := hex.EncodeToString(innerMd5[:])     // 转为字符串
	tokenSource := timestamp + innerMd5Str             // 拼接字符串
	token := md5.Sum([]byte(tokenSource))              // md5(timestamp + md5(api_sk))
	return hex.EncodeToString(token[:])                // 返回最终的 token
}

func (bt *BtApiClient) AppendSignatureBody(body url.Values) {
	now := time.Now()
	body.Set("request_time", fmt.Sprintf("%d", now.Unix()))
	body.Set("request_token", bt.GenerateRequestToken(now))
}

type RestartGoProjectResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"status_code"`
	ErrorMsg   string `json:"error_msg"`
	Data       string `json:"data"`
}

func (bt *BtApiClient) InvokeBtApi(path string, body url.Values, resp any) error {
	var err error
	bt.AppendSignatureBody(body)
	err = requests.
		URL(bt.btUrl).
		Transport(bt.ignoreSslTransport).
		Path(path).
		Method(http.MethodPost).
		BodyForm(body).
		ToJSON(resp).
		Fetch(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (bt *BtApiClient) RestartGoProject(name string) (*RestartGoProjectResponse, error) {
	var err error
	postForm := url.Values{}
	postForm.Set("project_name", name)
	var response RestartGoProjectResponse
	err = bt.InvokeBtApi("/project/go/restart_project", postForm, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type StartGoProjectResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"status_code"`
	ErrorMsg   string `json:"error_msg"`
	Data       string `json:"data"`
}

func (bt *BtApiClient) StartGoProject(name string) (*StartGoProjectResponse, error) {
	var err error
	postForm := url.Values{}
	postForm.Set("project_name", name)
	var response StartGoProjectResponse
	err = bt.InvokeBtApi("/project/go/start_project", postForm, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type StopGoProjectResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"status_code"`
	ErrorMsg   string `json:"error_msg"`
	Data       string `json:"data"`
}

func (bt *BtApiClient) StopGoProject(name string) (*StopGoProjectResponse, error) {
	var err error
	postForm := url.Values{}
	postForm.Set("project_name", name)
	var response StopGoProjectResponse
	err = bt.InvokeBtApi("/project/go/stop_project", postForm, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type StartNodeProjectResponse struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Status    bool   `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

func (bt *BtApiClient) StartNodeProject(name string) (*StartNodeProjectResponse, error) {
	var err error
	postForm := url.Values{}
	postForm.Set("project_name", name)
	postForm.Set("project_type", "general")
	postForm.Set("status", "start")
	var response StartNodeProjectResponse
	err = bt.InvokeBtApi("/mod/nodejs/com/set_project_status", postForm, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type StopNodeProjectResponse struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Status    bool   `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

func (bt *BtApiClient) StopNodeProject(name string) (*StopNodeProjectResponse, error) {
	var err error
	postForm := url.Values{}
	postForm.Set("project_name", name)
	postForm.Set("project_type", "general")
	postForm.Set("status", "stop")
	var response StopNodeProjectResponse
	err = bt.InvokeBtApi("/mod/nodejs/com/set_project_status", postForm, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type RestartNodeProjectResponse struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Status    bool   `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

func (bt *BtApiClient) RestartNodeProject(name string) (*RestartNodeProjectResponse, error) {
	var err error
	postForm := url.Values{}
	postForm.Set("project_name", name)
	postForm.Set("project_type", "general")
	postForm.Set("status", "restart")
	var response RestartNodeProjectResponse
	err = bt.InvokeBtApi("/mod/nodejs/com/set_project_status", postForm, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func NewBtApiClient(btUrl string, btKey string) *BtApiClient {
	return &BtApiClient{
		btUrl: btUrl,
		btKey: btKey,
		ignoreSslTransport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
