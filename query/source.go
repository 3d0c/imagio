package query

import (
	"bytes"
	"github.com/3d0c/imagio/config"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	NOPROTO = 1
	URIFULL = 2
)

var readers = map[string]func(string) (io.Reader, error){
	"http": http_reader,
	"file": file_reader,
}

type Source struct {
	BlobLen int

	scheme   string
	root     string
	filepath string
	blob     []byte
	reader   func(string) (io.Reader, error)
	Imgcfg   image.Config
	imgtype  string
}

func (*Source) Construct(i ...interface{}) *Source {
	if len(i) != 1 {
		log.Printf("Wrong arguments count = %d. Expecting 1\n", len(i))
		return nil
	}

	var err error
	v := i[0].([]interface{})[0].(string)

	this := &Source{}

	parts := strings.Split(v, "://")

	switch len(parts) {
	case NOPROTO:
		this.scheme = config.Get().Scheme()
		this.filepath = filepath.Clean(parts[0])
		break

	case URIFULL:
		this.scheme = parts[0]

		if strings.Contains(parts[1], "\x00") || strings.Contains(parts[1], "../") {
			log.Printf("URI containts illegal characters. `%v`\n", parts[1])
			return nil
		}

		this.filepath = filepath.Clean(parts[1])
		break

	default:
		log.Printf("Wrong URL = '%v'\n", v)
		return nil
	}

	this.root, err = config.Get().Root(this.scheme)
	if err != nil {
		log.Println("Unable to proceed without default root for `file` scheme.", err)
		return nil
	}

	if r, found := readers[this.scheme]; found {
		this.reader = r
	}

	if this.Blob() == nil {
		return nil
	}

	this.Imgcfg, this.imgtype, err = image.DecodeConfig(bytes.NewReader(this.Blob()))
	if err != nil {
		log.Printf("Unable to DecodeConfig() for resource from %v. %v\n", this.Link(), err)
		return nil
	}

	return this
}

func (this *Source) Blob() []byte {
	var err error

	if this.blob != nil {
		return this.blob
	}

	r, err := this.reader(this.Link())
	if err != nil {
		log.Printf("Unable to open resource for: %v. %v\n", this.Link(), err)
		return nil
	}

	this.blob, err = ioutil.ReadAll(r)
	if err != nil {
		log.Printf("Unable to read resource for: %v. %v\n", this.Link(), err)
		return nil
	}

	this.BlobLen = len(this.blob)

	return this.blob
}

func (this *Source) Config() *image.Config {
	return &this.Imgcfg
}

func (this *Source) Type() string {
	return this.imgtype
}

func (this *Source) Mime() string {
	return mime.TypeByExtension("." + this.Type())
}

func (this *Source) LinkFull() string {
	return this.scheme + "://" + this.root + this.filepath
}

func (this *Source) Link() string {
	return this.root + this.filepath
}

func (this *Source) Size() *PixelDim {
	return &PixelDim{Width: this.Imgcfg.Width, Height: this.Imgcfg.Height}
}

func http_reader(src string) (io.Reader, error) {
	res, err := http.Get("http://" + src)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func file_reader(src string) (io.Reader, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}

	return file, nil
}
