package subscriber

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/micro/go-micro/v2/codec/proto"
	"github.com/micro/go-micro/v2/config"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"

	ufile "github.com/bimoyong/go-util/file"
	umetadata "github.com/bimoyong/go-util/metadata"
)

// File struct
type File struct {
}

// OnMessage function
func (s *File) OnMessage(ctx context.Context, message *proto.Message) (err error) {
	md, _ := metadata.FromContext(ctx)
	id, _ := md.Get("ID")
	domain, _ := md.Get("Domain")
	defer func() {
		if err != nil {
			log.Errorf("[File] Store failed!: err=[%s] metadata=[%+v] data=[%s]", err.Error(), md, message.Data)
		}
	}()

	var msg map[string]string
	if err = json.Unmarshal(message.Data, &msg); err != nil {
		err = fmt.Errorf("Cannot decode message: err=[%s]", err.Error())

		return
	}
	var field string
	var data string
	for field, data = range msg {
		break
	}
	if field == "" {
		return
	}

	var fileName string
	dirBase := filepath.Join(config.Get("dir_base").String(""), domain)
	if fileName, err = generateFileName(data, dirBase, md); err != nil {
		err = fmt.Errorf("Cannot generate file name: err=[%s] dir_base=[%s]", err.Error(), dirBase)

		return
	}
	log.Debug("[File] Generate file name: ", fileName)

	r := strings.NewReader(data)
	if _, err = save2Disk(fileName, r); err != nil {
		err = fmt.Errorf("Cannot save to disk: err=[%s] file_name=[%s]", err.Error(), fileName)

		return
	}
	log.Infof("[File] Save to disk success: id=[%s], file_name=[%s]", id, fileName)

	postback := map[string]interface{}{
		field: fileName,
	}
	md.Set("Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	if err = publish(postback, umetadata.NewMetadata(md)); err != nil {
		log.Errorf("[File] Postback failed!: metadata=[%+v] msg=[%+v]", md, postback)
	} else {
		log.Debugf("[File] Postback: metadata=[%+v] msg=[%+v]", md, postback)
	}

	return
}

func generateFileName(data string, base string, md metadata.Metadata) (fileName string, err error) {
	domain, _ := md.Get("Domain")
	alias, _ := md.Get("Alias")
	resource, _ := md.Get("Resource")
	base = filepath.Join(base, domain, alias, resource, time.Now().Format("2006-01-02"))
	if err = ufile.CheckOrMkdirAll(base); err != nil {
		return
	}

	var b []byte
	b, err = base64.StdEncoding.DecodeString(data)
	r := bytes.NewReader(b)
	var ext string
	if _, ext, err = image.DecodeConfig(r); err != nil {
		return
	}

	fileName = filepath.Join(base, hash(b)+"."+ext)

	return
}

func save2Disk(fileName string, reader io.Reader) (size int64, err error) {
	var f *os.File
	if f, err = os.Create(fileName); err != nil {
		return
	}
	f.SetWriteDeadline(time.Now().Add(time.Second * 60))
	defer f.Close()

	dec := base64.NewDecoder(base64.StdEncoding, reader)
	size, err = io.Copy(f, dec)

	return
}

func hash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}
