package main

import (
	"context"
	"log"
	"os"

	"github.com/hempflower/btcli/internal/btapi"
	"github.com/urfave/cli/v3"
)

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
						Name: "start",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "The name of the go project",
								Required: true,
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							bt := btapi.NewBtApiClient(cmd.String("bt-url"), cmd.String("bt-key"))
							response, err := bt.StartGoProject(cmd.String("name"))
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
					{
						Name: "stop",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "The name of the go project",
								Required: true,
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							bt := btapi.NewBtApiClient(cmd.String("bt-url"), cmd.String("bt-key"))
							response, err := bt.StopGoProject(cmd.String("name"))
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
							bt := btapi.NewBtApiClient(cmd.String("bt-url"), cmd.String("bt-key"))
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
			{
				Name: "node",
				Commands: []*cli.Command{
					{
						Name: "start",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "The name of the go project",
								Required: true,
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							bt := btapi.NewBtApiClient(cmd.String("bt-url"), cmd.String("bt-key"))
							response, err := bt.StartNodeProject(cmd.String("name"))
							if err != nil {
								return err
							}
							if !response.Status {
								log.Fatal(response.Msg)
							} else {
								log.Println(response.Msg)
							}
							return nil
						},
					},
					{
						Name: "stop",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "The name of the go project",
								Required: true,
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							bt := btapi.NewBtApiClient(cmd.String("bt-url"), cmd.String("bt-key"))
							response, err := bt.StopNodeProject(cmd.String("name"))
							if err != nil {
								return err
							}
							if !response.Status {
								log.Fatal(response.Msg)
							} else {
								log.Println(response.Msg)
							}
							return nil
						},
					},
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
							bt := btapi.NewBtApiClient(cmd.String("bt-url"), cmd.String("bt-key"))
							response, err := bt.RestartNodeProject(cmd.String("name"))
							if err != nil {
								return err
							}
							if !response.Status {
								log.Fatal(response.Msg)
							} else {
								log.Println(response.Msg)
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
