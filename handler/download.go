package handler

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/micro/go-micro/v2/config"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/bimoyong/go-file/proto/file"
)

// Download function
func (h *File) Download(ctx context.Context, req *proto.DownloadReq, strm proto.File_DownloadStream) (err error) {
	log.Info("Download done!")

	md, _ := metadata.FromContext(ctx)

	// check chunk size given by client is not over server's limit
	var chunk_size_int int64
	chunk_size_limit := int64(config.Get("chunk_size_limit").Int(1 << 20))
	chunk_size, ok := md.Get("chunk_size")
	if ok {
		chunk_size_int, err = strconv.ParseInt(chunk_size, 10, 64)
		if err != nil || chunk_size_int <= 0 || chunk_size_limit < chunk_size_int {
			ok = false
		}
	}
	if !ok {
		chunk_size_int = chunk_size_limit
	}

	defer func() {
		if err != nil {
			log.Errorf("Failed to send file! err=[%s] metadata=[%+v] fileinfo=[%+v]", err.Error(), md, "fileinfo")

			return
		}
	}()

	var resp proto.DownloadResp
	var size = int64(0)
	buf := make([]byte, chunk_size_int)
	f, _ := os.Open("../resource_test/12032006.jpeg")

	finfo, err := f.Stat()
	if err != nil {
		log.Errorf("Error stating file! %s", err.Error())
	}

	log.Infof("Send file for every chunk size %d", cap(buf))
	for {
		var n int
		n, err = f.Read(buf)
		if err == io.EOF {
			err = nil
			log.Infof("Finished sending file")
			break
		}
		if err != nil {
			err = fmt.Errorf("error reading file: %s", err.Error())
			break
		}

		size += int64(n)
		if size >= finfo.Size() {
			resp.Desc = &proto.Description{
				Ext:       filepath.Ext(finfo.Name()),
				Size:      finfo.Size(),
				CreatedAt: timestamppb.New(finfo.ModTime()),
			}
		}
		resp.Data = buf[:n]
		resp.Checksum = fmt.Sprintf("%x", sha1.Sum(buf[:n]))
		resp.Timestamp = timestamppb.New(time.Now())

		if err = strm.SendMsg(&resp); err != nil {
			err = fmt.Errorf("error sending file: %s", err.Error())
			break
		}
	}

	return
}
