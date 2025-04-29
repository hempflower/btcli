package main

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/urfave/cli/v3"
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

func (bt *BtApiClient) RestartGoProject(name string) (*RestartGoProjectResponse, error) {
	var err error
	postForm := url.Values{}
	postForm.Set("project_name", name)
	bt.AppendSignatureBody(postForm)
	var response RestartGoProjectResponse
	err = requests.
		URL(bt.btUrl).
		Transport(bt.ignoreSslTransport).
		Path("/project/go/restart_project").
		Method("POST").
		BodyForm(postForm).
		ToJSON(&response).
		Fetch(context.Background())
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

func main() {
	cmd := &cli.Command{
		Name:  "btcli",
		Usage: "A command-line tool for managing BT-Panel",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bt-url",
				Usage:    "The URL of the BT-Panel API",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "bt-key",
				Usage:    "The API key for the BT-Panel API",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name: "go-project",
				Commands: []*cli.Command{
					{
						Name: "restart",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "The name of the go project",
								Required: true,
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							bt := NewBtApiClient(cmd.String("bt-url"), cmd.String("bt-key"))
							response, err := bt.RestartGoProject(cmd.String("name"))
							if err != nil {
								return err
							}
							if !response.Status {
								log.Fatal(response.ErrorMsg)
							} else {
								log.Println(response.Data)
							}
							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
