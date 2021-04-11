package handler

import (
	"context"

	proto "github.com/bimoyong/go-file/proto/file"
)

// Download function
func (h *File) Download(ctx context.Context, req *proto.DownloadReq, stream proto.File_DownloadStream) (err error) {
	// md, _ := metadata.FromContext(ctx)
	// var resp proto.UploadResp
	// var name string
	// var file *os.File
	// var fileinfo os.FileInfo
	// var size int
	// var sizeMax = config.Get("bytes_limit").Int(5 << 20)

	// defer func() {
	// 	if err != nil && err != io.EOF {
	// 		log.Errorf("Failed to receive file! err=[%s] metadata=[%+v] fileinfo=[%+v]", err.Error(), md, fileinfo)

	// 		_ = os.RemoveAll(name)

	// 		return
	// 	}

	// }()

	// for {
	// 	var chunk *proto.UploadReq
	// 	chunk, err = stream.Recv()
	// 	if err == io.EOF {
	// 		log.Infof("Finished receiving file %s", name)
	// 		break
	// 	}
	// 	if err != nil {
	// 		err = status.Errorf(codes.Unknown, err.Error())
	// 		return
	// 	}

	// 	if checksum := fmt.Sprintf("%x", sha1.Sum(chunk.Data)); checksum != chunk.Checksum {
	// 		err = status.Errorf(codes.DataLoss, "incorrect checksum: expect %s but given %s", chunk.Checksum, checksum)
	// 		return
	// 	}

	// 	if size += len(chunk.Data); size > sizeMax {
	// 		err = status.Errorf(codes.ResourceExhausted, "file is too large: %d > %d", size, sizeMax)
	// 		return
	// 	}
	// 	log.Debugf("Received %d bytes of file %s", size, name)

	// 	if file == nil {
	// 		base := filepath.Join(config.Get("dir_base").String(""), md["Alias"])
	// 		if name, err = util.NewName(chunk.Data, base); err != nil {
	// 			err = status.Errorf(codes.Internal, "error determining file name: %s", err.Error())
	// 			return
	// 		}
	// 		log.Debugf("Generate file name %s", name)

	// 		if file, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755); err != nil {
	// 			err = status.Errorf(codes.Internal, "error opening file %s: %s", name, err.Error())
	// 			return
	// 		}
	// 		defer file.Close()

	// 		fileinfo, _ = file.Stat()
	// 		tm, _ := ptypes.TimestampProto(fileinfo.ModTime())
	// 		resp = proto.UploadResp{
	// 			Id:        strings.TrimSuffix(fileinfo.Name(), path.Ext(fileinfo.Name())),
	// 			Timestamp: tm,
	// 		}
	// 	}

	// 	if _, err = file.Write(chunk.Data); err != nil {
	// 		err = status.Errorf(codes.Internal, "error writing file name %s: %s", name, err.Error())
	// 		return
	// 	}

	// 	if err = stream.SendMsg(&resp); err != nil {
	// 		return
	// 	}
	// }

	return
}
