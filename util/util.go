package util

import (
	"math"
	"mime"
	"net/http"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson/primitive"

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
