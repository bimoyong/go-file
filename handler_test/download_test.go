package handler_test

import (
	"context"
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
	"github.com/bimoyong/go-file/util"
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
	ctx = metadata.Set(ctx, "Chunk-Size", "102400")
	md, _ := metadata.FromContext(ctx)

	req := proto.DownloadReq{
		Id: "6131f2f9ccb6e0ba045e07b4",
	}
	srv := proto.NewFileService(test.ServerName, client.DefaultClient)
	strm, err := srv.Download(ctx, &req)
	if err != nil {
		t.Fatalf("Error handshaking download service! %s", err.Error())
	}
	defer strm.Close()
	t.Logf("Setup! metadata=[%+v] req=[%s]", md, req.String())

	fname := filepath.Join(os.TempDir(), fmt.Sprintf("%s.jpg", req.Id))
	f, err := os.Create(fname)
	if err != nil {
		t.Fatalf("Error creating file %s! %s", fname, err.Error())
	}
	defer f.Close()

	var resp proto.DownloadResp
	size := int64(0)

	t.Logf("Start downloading file %s", fname)
	for {
		err := strm.RecvMsg(&resp)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error downloading! %s", err.Error())
		}

		n, err := f.Write(resp.Chunk.Data)
		if err != nil {
			t.Fatalf("Error writing file! %s", err.Error())
		}

		if checksum, err := util.Checksum(resp.Chunk.Checksum, resp.Chunk.Data); err != nil {
			t.Fatalf("Incorrect checksum! Expect %s but given %s", resp.Chunk.Checksum, checksum)
		}

		size += int64(n)
		t.Logf("Downloaded %d bytes of file", size)
	}

	t.Logf("File info: %s", resp.Desc)

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
