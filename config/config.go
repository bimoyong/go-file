package config

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/config/reader"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"

	"github.com/bimoyong/go-util/config"
	ufile "github.com/bimoyong/go-util/file"
)

// Config struct
type Config struct{}

// Init function
func Init() (err error) {
	if err = newBroker(); err != nil {
		return
	}

	config.Init(&Config{})

	return
}

// DirBase function
func (s *Config) DirBase(value reader.Value) {
	v := value.String("")
	if len(v) > 0 {
		if err := ufile.CheckOrMkdirAll(v); err != nil {
			log.Error("[Config][DirBase] Create directory failed!: ", v)
		}
	}
}

func newBroker() error {
	broker.DefaultBroker = server.DefaultServer.Options().Broker

	return broker.DefaultBroker.Connect()
}
