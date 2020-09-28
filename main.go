package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"

	"gitlab.com/bimoyong/go-file/config"
	sub "gitlab.com/bimoyong/go-file/subscriber"
)

func main() {
	service := micro.NewService()

	client.DefaultClient = service.Client()
	server.DefaultServer = service.Server()

	service.Init(
		micro.BeforeStart(config.Init),
		micro.BeforeStart(sub.RegisterFile),
		micro.BeforeStop(sub.Close),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
