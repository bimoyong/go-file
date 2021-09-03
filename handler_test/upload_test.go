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
	"github.com/micro/go-micro/v2/metadata"

	proto "github.com/bimoyong/go-file/proto/file"
	test "github.com/bimoyong/go-file/test"
)

// TestUpload function
func TestUpload(t *testing.T) {
	go func() {
		if err := test.StartService(); err != nil {
			t.Errorf("Error starting service! %s", err.Error())
		}
	}()
	time.Sleep(time.Second * 2)

	ctx := metadata.Set(context.TODO(), "Domain", "staging")
	ctx = metadata.Set(ctx, "Alias", "vehicles")
	md, _ := metadata.FromContext(ctx)

	req := proto.UploadReq{
		Checksum: "todo_checksum",
	}
	t.Logf("Setup! metadata=[%+v] req=[%s]", md, req.String())

	srv := proto.NewFileService(test.ServerName, client.DefaultClient)
	strm, err := srv.Upload(ctx)
	if err != nil {
		t.Fatalf("Error handshaking upload service! %s", err.Error())
	}
	defer strm.Close()

	fname := "../resource_test/12032006.jpeg"
	f, err := os.Open(fname)
	if err != nil {
		t.Fatalf("Error opening file %s! %s", fname, err.Error())
	}
	defer f.Close()

	buf := make([]byte, 1<<20)
	size := int64(0)
	resp := proto.UploadResp{}

	t.Logf("Start uploading file for every chunk size %d", cap(buf))
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error reading file! %s", err.Error())
		}

		req.Data = buf[:n]
		req.Checksum = fmt.Sprintf("%x", sha1.Sum(req.Data))
		if err = strm.Send(&req); err != nil {
			t.Fatalf("Error uploading! %s", err.Error())
		}

		size += int64(n)
		t.Logf("Uploaded %d bytes", size)

		if err = strm.RecvMsg(&resp); err != nil {
			t.Fatalf("Error ack uploading! %s", err.Error())
		}
	}

	t.Logf("Server responds %s", resp.String())

	finfo, err := f.Stat()
	if err != nil {
		t.Fatalf("Error stating file! %s", err.Error())
	}
	if finfo.Size() != size {
		t.Fatalf("File is not uploaded fully! Expect %d bytes", finfo.Size())
	}

	t.Logf("Teardown!")
}
