package subscriber

import (
	"bytes"
	"context"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/codec/proto"
	"github.com/micro/go-micro/v2/config"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"

	ufile "github.com/bimoyong/go-util/file"
	umetadata "github.com/bimoyong/go-util/metadata"
	"gitlab.com/bimoyong/go-file/model"
)

// File struct
type File struct {
}

// OnMessage function
func (s *File) OnMessage(ctx context.Context, message *proto.Message) (err error) {
	md, _ := metadata.FromContext(ctx)
	defer func() {
		if err != nil {
			log.Errorf("[File] Store failed!: err=[%s] metadata=[%+v] data=[%s]", err.Error(), md, message.Data)
		}
	}()

	var msg model.Message
	if err = json.Unmarshal(message.Data, &msg); err != nil {
		err = fmt.Errorf("Cannot decode message: err=[%s]", err.Error())

		return
	}

	id, _ := md.Get("ID")
	domain, _ := md.Get("Domain")
	postback, _ := md.Get("Postback")

	var fileName string
	dirBase := filepath.Join(config.Get("dir_base").String(""), domain)
	if fileName, err = generateFileName(msg, dirBase); err != nil {
		err = fmt.Errorf("Cannot generate file name: err=[%s] dir_base=[%s]", err.Error(), dirBase)

		return
	}
	log.Debug("[File] Generate file name: ", fileName)

	if err = save2Disk(msg.Kind, fileName, msg.File); err != nil {
		err = fmt.Errorf("Cannot save to disk: err=[%s] file_name=[%s]", err.Error(), fileName)

		return
	}
	log.Infof("[File] Save to disk success: id=[%s], file_name=[%s]", id, fileName)

	var pub micro.Publisher
	pb := model.Postback{
		// ID:        h.ID,
		Name:      fileName[len(config.Get("dir_base").String("")):],
		FullName:  fileName,
		Timestamp: time.Now(),
	}
	if pub, err = postback1(pb, umetadata.NewMetadata(md)); err != nil {
		log.Errorf("[File] Postback failed!: metadata=[%+v] msg=[%+v]", md, pb)
	} else {
		log.Debugf("[File] Postback: metadata=[%+v] msg=[%+v]", md, pb)
	}
	if _, ok := PostbackMap[postback]; !ok {
		PostbackMap[postback] = pub
	}

	return
}

func generateFileName(m model.Message, base string) (fileName string, err error) {
	var b []byte
	if m.Name == "" {
		b, err = base64.StdEncoding.DecodeString(m.File)
		m.Name = hash(b)
	}

	mid, name := filepath.Split(m.Name)
	fileName = filepath.Join(base, mid, time.Now().Format("2006-01-02"), name)
	base, _ = filepath.Split(fileName)
	if err = ufile.CheckOrMkdirAll(base); err != nil {
		return
	}

	// TODO: make this more flexible
	var ext string
	switch {
	case m.Kind.Enabled(model.Base64Kind):
		r := bytes.NewReader(b)
		if _, ext, err = image.DecodeConfig(r); err != nil {
			return
		}
		fileName = fileName + "." + ext
	case m.Kind.Enabled(model.URLKind):
		if "" == filepath.Ext(fileName) {
			fileName = fileName + filepath.Ext(m.File)
		}
	}

	return
}

func save2Disk(kind model.Kind, fileName string, data string) (err error) {
	// TODO: make this more flexible
	if kind.Enabled(model.URLKind) {
		err = download2Disk(fileName, data)

		return
	}

	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))

	_, err = writeFile(fileName, dec)

	return
}

func download2Disk(fileName string, data string) (err error) {
	cli := http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	var rsp *http.Response
	if rsp, err = cli.Get(data); err != nil {
		return
	}
	defer rsp.Body.Close()

	_, err = writeFile(fileName, rsp.Body)

	return
}

func writeFile(name string, reader io.Reader) (size int64, err error) {
	var f *os.File
	if f, err = os.Create(name); err != nil {
		return
	}
	f.SetWriteDeadline(time.Now().Add(time.Second * 60))
	defer f.Close()

	size, err = io.Copy(f, reader)

	return
}

func hash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}
