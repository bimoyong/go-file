package handler

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/micro/go-micro/v2/config"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/bimoyong/go-file/proto/file"
	"github.com/bimoyong/go-file/util"
)

// Download function
func (h *File) Download(ctx context.Context, req *proto.DownloadReq, strm proto.File_DownloadStream) (err error) {
	md, _ := metadata.FromContext(ctx)

	defer func() {
		if err != nil {
			log.Errorf("Failed to send file! err=[%s] metadata=[%+v] fileinfo=[%+v]", err.Error(), md, req.Id)

			return
		}
	}()

	resp := proto.DownloadResp{Chunk: &proto.Chunk{}}
	var size = int64(0)
	buf := make([]byte, util.DetermineChunkSize(md))

	var id primitive.ObjectID
	if id, err = primitive.ObjectIDFromHex(req.Id); err != nil {
		log.Warnf("Magic request! %s", err.Error())
		err = status.Errorf(codes.NotFound, "ID %s not found!", req.Id)
		return
	}
	id.Timestamp().Format("2006-01-02")
	fname := filepath.Join(
		config.Get("dir_base").String(""),
		md["Domain"], md["Alias"],
		id.Timestamp().Format("2006-01-02"),
		req.Id+".jpeg",
	)
	var finfo os.FileInfo
	if finfo, err = os.Stat(fname); os.IsNotExist(err) {
		log.Errorf("File does not exist! %s", err.Error())
		err = status.Errorf(codes.NotFound, "ID %s not found!", req.Id)
		return
	}
	f, _ := os.Open(fname)
	defer f.Close()

	log.Infof("Send file for every chunk size %d", cap(buf))
	for {
		var n int
		if n, err = f.Read(buf); err == io.EOF {
			log.Infof("Finished sending file")
			err = nil
			break
		}
		if err != nil {
			log.Errorf("Error reading file! %s", err.Error())
			err = status.Error(codes.Internal, "Internal server error")
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
		resp.Chunk.Data = buf[:n]
		resp.Chunk.Checksum = fmt.Sprintf("%x", sha1.Sum(buf[:n]))
		resp.Timestamp = timestamppb.New(time.Now())

		if err = strm.SendMsg(&resp); err != nil {
			log.Errorf("Error sending file %s", err.Error())
			break
		}
	}

	return
}
