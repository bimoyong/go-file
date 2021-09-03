package test

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/server"

	"github.com/bimoyong/go-file/handler"
	proto "github.com/bimoyong/go-file/proto/file"
	userver "github.com/bimoyong/go-util/server"
)

const ServerName = "go.srv.file"

func StartService() error {
	config.DefaultConfig.Set("./data", "dir_base")
	config.DefaultConfig.Set(500<<20, "bytes_limit")

	service := micro.NewService()

	server.DefaultServer = service.Server()
	client.DefaultClient = service.Client()

	server.DefaultServer.Init(
		server.WrapSubscriber(userver.LogWrapper),
		server.WrapSubscriber(userver.AuthWrapper),
	)
	server.DefaultServer.Init(
		server.Name(ServerName),
	)

	proto.RegisterFileHandler(server.DefaultServer, new(handler.File))

	return service.Run()
}
