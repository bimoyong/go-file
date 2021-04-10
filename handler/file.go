package handler

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/ptypes"
	"github.com/micro/go-micro/v2/config"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "github.com/bimoyong/go-file/proto/file"
	"github.com/bimoyong/go-file/util"
)

// File struct
type File struct {
}

// Upload function
func (h *File) Upload(ctx context.Context, stream proto.File_UploadStream) (err error) {
	md, _ := metadata.FromContext(ctx)
	var name string
	var file *os.File
	var fileinfo os.FileInfo
	var size int
	var sizeMax = config.Get("bytes_limit").Int(5 << 20)

	defer func() {
		if err != nil {
			log.Errorf("Failed to receive file! err=[%s] metadata=[%+v] fileinfo=[%+v]", err.Error(), md, fileinfo)
		} else {
			log.Infof("Finished receiving file %s", name)
		}
	}()

	for {
		var chunk *proto.UploadReq
		chunk, err = stream.Recv()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			err = status.Errorf(codes.Unknown, err.Error())
			return
		}

		if size += len(chunk.Data); size > sizeMax {
			err = status.Errorf(codes.ResourceExhausted, "file is too large: %d > %d", size, sizeMax)
			return
		}
		log.Debugf("Received %d bytes of file %s", size, name)

		if file == nil {
			base := filepath.Join(config.Get("dir_base").String(""), md["Alias"])
			if name, err = util.NewName(chunk.Data, base); err != nil {
				err = fmt.Errorf("error determining file name: %s", err.Error())
				return
			}
			log.Debugf("Generate file name %s", name)

			if file, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755); err != nil {
				err = fmt.Errorf("error opening file %s: %s", name, err.Error())
				return
			}

			defer file.Close()
		}

		if _, err = file.Write(chunk.Data); err != nil {
			err = fmt.Errorf("error writing file name %s: %s", name, err.Error())
			return
		}

		fileinfo, _ = file.Stat()
		tm, _ := ptypes.TimestampProto(fileinfo.ModTime())
		resp := proto.UploadResp{
			Id:        fileinfo.Name(),
			Timestamp: tm,
		}
		if err = stream.SendMsg(&resp); err != nil {
			return
		}
	}

	return
}
