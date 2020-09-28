package subscriber

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/config"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"

	raw "github.com/micro/go-micro/v2/codec/bytes"
	"gitlab.com/bimoyong/go-file/model"
)

var (
	// PostbackMap stores created publishers
	PostbackMap = map[string]micro.Publisher{}
)

// RegisterFile function
func RegisterFile() error {
	return micro.RegisterSubscriber(
		config.Get("broker", "topic_in").String(""),
		server.DefaultServer,
		&File{},
		server.SubscriberQueue(config.Get("broker", "queue").String("")),
	)
}

// Close function
func Close() error {
	log.Info("[Subscriber][Close] Do nothing")

	return nil
}

func postback1(m model.Postback, md metadata.Metadata) (pub micro.Publisher, err error) {
	topic, ok := md.Get("Postback")
	if !ok {
		return
	}
	md.Delete("Postback")

	p := micro.NewPublisher(topic, client.DefaultClient)
	ctx := metadata.NewContext(context.Background(), md)
	msg := raw.Frame{}
	msg.Data, _ = json.Marshal(m)
	if err = p.Publish(ctx, &msg); err != nil {
		return
	}
	pub = p

	return
}
