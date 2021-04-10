package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"

	"github.com/bimoyong/go-file/config"
	"github.com/bimoyong/go-file/handler"
	proto "github.com/bimoyong/go-file/proto/file"
	sub "github.com/bimoyong/go-file/subscriber"
	userver "github.com/bimoyong/go-util/server"
)

func main() {
	service := micro.NewService()

	client.DefaultClient = service.Client()
	server.DefaultServer = service.Server()

	server.DefaultServer.Init(
		server.WrapSubscriber(userver.LogWrapper),
		server.WrapSubscriber(userver.AuthWrapper),
	)

	service.Init(
		micro.BeforeStart(config.Init),
		micro.BeforeStart(sub.RegisterFile),
		micro.BeforeStop(sub.Close),
	)

	proto.RegisterFileHandler(server.DefaultServer, new(handler.File))

	if err := service.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
