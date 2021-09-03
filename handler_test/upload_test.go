package handler_test

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"

	proto "github.com/bimoyong/go-file/proto/file"
	test "github.com/bimoyong/go-file/test"
)

// TestUpload function
func TestUpload(t *testing.T) {
	var err error

	go func() {
		if err = test.StartService(); err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second * 2)

	ctx := metadata.Set(context.TODO(), "Domain", "staging")
	ctx = metadata.Set(ctx, "Alias", "vehicles")

	f, _ := os.Open("../resource_test/12032006.jpeg")
	req := proto.UploadReq{
		Checksum: "todo_checksum",
	}

	if err = upload(ctx, &req, f); err != nil {
		t.Error(err)
	}
	t.Logf("Test done")
}

func upload(ctx context.Context, req *proto.UploadReq, rd io.Reader) (err error) {
	srv := proto.NewFileService(test.ServerName, client.DefaultClient)
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
