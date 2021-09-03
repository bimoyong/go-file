package handler_test

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"

	proto "github.com/bimoyong/go-file/proto/file"
	test "github.com/bimoyong/go-file/test"
)

// TestDownload function
func TestDownload(t *testing.T) {
	go func() {
		if err := test.StartService(); err != nil {
			t.Errorf("Error starting service! %s", err.Error())
		}
	}()
	time.Sleep(time.Second * 2)

	ctx := metadata.Set(context.TODO(), "Domain", "staging")
	ctx = metadata.Set(ctx, "Alias", "vehicles")
	md, _ := metadata.FromContext(ctx)

	req := proto.DownloadReq{
		Id: "todo_id",
	}
	t.Logf("Setup! metadata=[%+v] req=[%s]", md, req.String())

	srv := proto.NewFileService(test.ServerName, client.DefaultClient)
	strm, err := srv.Download(ctx, &req)
	if err != nil {
		t.Fatalf("Error handshaking download service! %s", err.Error())
	}
	defer strm.Close()

	fname := filepath.Join(os.TempDir(), fmt.Sprintf("%s.jpg", req.Id))
	f, err := os.Create(fname)
	if err != nil {
		t.Fatalf("Error creating file %s! %s", fname, err.Error())
	}
	defer f.Close()

	size := int64(0)
	resp := proto.DownloadResp{}

	t.Logf("Start downloading file %s", fname)
	for {
		err := strm.RecvMsg(&resp)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error downloading! %s", err.Error())
		}

		n, err := f.Write(resp.Data)
		if err != nil {
			t.Fatalf("Error writing file! %s", err.Error())
		}
		checksum := fmt.Sprintf("%x", sha1.Sum(resp.Data))
		if checksum != resp.Checksum {
			t.Fatalf("Incorrect checksum! Expect %s but given %s", checksum, resp.Checksum)
		}

		size += int64(n)
		t.Logf("Downloaded %d bytes of file", size)
	}

	finfo, err := f.Stat()
	if err != nil {
		t.Fatalf("Error stating file! %s", err.Error())
	}
	if finfo.Size() != 577855 {
		t.Errorf("File is not downloaded fully! Expect %d bytes", 577855)
	}

	err = os.Remove(fname)
	if err != nil {
		t.Fatalf("Error delete file! %s", err.Error())
	}
	t.Logf("Teardown!")
}
