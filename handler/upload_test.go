package handler

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/config"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"

	"github.com/bimoyong/go-file/handler"
	proto "github.com/bimoyong/go-file/proto/file"
)

// TestUpload function
func TestUpload(t *testing.T) {
	var err error

	go func() {
		if err = startService(); err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second)

	ctx := metadata.Set(context.TODO(), "Alias", "vehicles")

	f, _ := os.Open("./upload_test.jpeg")
	req := proto.UploadReq{
		Checksum: "todo_checksum",
	}

	if err = upload(ctx, &req, f); err != nil {
		t.Error(err)
	}
	t.Logf("Test done")

	return
}

func upload(ctx context.Context, req *proto.UploadReq, rd io.Reader) (err error) {
	srv := proto.NewFileService("go.srv.file", client.DefaultClient)
	var strm proto.File_UploadService
	if strm, err = srv.Upload(ctx); err != nil {
		err = fmt.Errorf("failed to upload file: %s", err.Error())
		return
	}
	defer strm.Close()

	var resp proto.UploadResp

	// sent data in 1M chunks
	buf := make([]byte, 1<<20)

	log.Infof("Stream file for every chunk size %d", 1<<20)
	for {
		var n int
		if n, err = rd.Read(buf); err == io.EOF {
			err = nil
			log.Infof("Finished uploading file")
			break
		}
		if err != nil {
			err = fmt.Errorf("error reading file: %s", err.Error())
			break
		}

		req.Data = buf[:n]
		req.Checksum = fmt.Sprintf("%x", sha1.Sum(buf[:n]))
		if err = strm.Send(req); err != nil {
			err = fmt.Errorf("error uploading file: %s", err.Error())
			break
		}
		log.Infof("Streamed %d bytes of file", n)

		if err = strm.RecvMsg(&resp); err != nil {
			return
		}
	}
	log.Infof("Received response %s", resp.String())

	return
}

func startService() error {
	config.DefaultConfig.Set("../data", "dir_base")
	config.DefaultConfig.Set(524288000, "bytes_limit")

	service := micro.NewService()

	server.DefaultServer = service.Server()
	server.DefaultServer.Init(
		server.Name("go.srv.file"),
	)

	proto.RegisterFileHandler(server.DefaultServer, new(handler.File))

	return service.Run()
}
