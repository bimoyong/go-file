package util

import (
	"crypto/sha1"
	"fmt"
	"math"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/metadata"
	"go.mongodb.org/mongo-driver/bson/primitive"

	proto "github.com/bimoyong/go-file/proto/file"
	ufile "github.com/bimoyong/go-util/file"
)

// NewName function
func NewName(base string, buffer ...[]byte) (name string, err error) {
	id := primitive.NewObjectID()
	sub := id.Timestamp().Format("2006-01-02")

	dir := filepath.Join(base, sub)
	if err = ufile.CheckOrMkdirAll(dir); err != nil {
		return
	}

	var ext string
	if len(buffer) > 0 {
		if ext, err = DetectExtension(buffer[0]); err != nil {
			return
		}
	}

	name = filepath.Join(dir, id.Hex()) + ext

	return
}

// DetectExtension function
func DetectExtension(buffer []byte) (extension string, err error) {
	// Only the first 512 bytes are used to sniff the content type.
	l := math.Min(float64(len(buffer)), 512)
	cntType := http.DetectContentType(buffer[:int(l)])

	var exts []string
	if exts, err = mime.ExtensionsByType(cntType); err != nil {
		return
	}

	if len(exts) > 0 {
		extension = exts[0]
	}

	return
}

// DetermineChunkSize function return chunk size given by client but not over server's limit
func DetermineChunkSize(md metadata.Metadata) (chunk_size_int int64) {
	var err error

	chunk_size_limit := int64(config.Get("chunk_size_limit").Int(1 << 20))
	chunk_size, ok := md.Get("Chunk-Size")
	if ok {
		chunk_size_int, err = strconv.ParseInt(chunk_size, 10, 64)
		if err != nil || chunk_size_int <= 0 || chunk_size_limit < chunk_size_int {
			ok = false
		}
	}
	if !ok {
		chunk_size_int = chunk_size_limit
	}

	return
}

func Checksum(chunk *proto.Chunk, data []byte) (err error) {
	checksum := fmt.Sprintf("%x", sha1.Sum(chunk.Data))
	if checksum != chunk.Checksum {
		err = fmt.Errorf("expect %s but given %s", checksum, chunk.Checksum)
	}

	return
}
