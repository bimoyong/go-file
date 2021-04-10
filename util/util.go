package util

import (
	"mime"
	"net/http"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson/primitive"

	ufile "github.com/bimoyong/go-util/file"
)

// NewName function
func NewName(buffer []byte, base string) (name string, err error) {
	id := primitive.NewObjectID()
	sub := id.Timestamp().Format("2006-01-02")

	dir := filepath.Join(base, sub)
	if err = ufile.CheckOrMkdirAll(dir); err != nil {
		return
	}

	// Only the first 512 bytes are used to sniff the content type.
	cntType := http.DetectContentType(buffer[:512])

	var exts []string
	if exts, err = mime.ExtensionsByType(cntType); err != nil {
		return
	}

	name = filepath.Join(dir, id.Hex())
	if len(exts) > 0 {
		name += exts[0]
	}

	return
}
