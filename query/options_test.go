package query

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"github.com/3d0c/imagio/config"
	. "github.com/3d0c/imagio/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

const test_server = "localhost:41899"
const file_name = "1024x768.jpg"
const sample_jpeg = "1f8b080866f18152020331303234783736382e6a706700edce3b0ec2301045d13718210a8a588226280d8515576e421b2b418a050bcd3a2858049f869d1807d1f069d2a277dc5dc93313cff18eecd0ed3b8800921ee20d3bacf3dc16b631a6699d736df075edc308c31029acad4ce5cbd207bf1df5fd35e4083d579842c906132d4a4b3c61f53cf5cd2ce565f6593154f959f577bd62a1246d511a1e3d8888888888888888e89f49bc3c005c75237d1b130000"

func getJpeg() []byte {
	img_gz, err := hex.DecodeString(sample_jpeg)
	if err != nil {
		log.Fatal("Unable to read sample_jpeg.", err)
	}

	r, err := gzip.NewReader(bytes.NewReader(img_gz))
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func serveJpeg() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, file_name, time.Now(), bytes.NewReader(getJpeg()))
	})

	log.Fatal(http.ListenAndServe(test_server, nil))
}

func init() {
	go serveJpeg()

	tmp, err := os.Create("/tmp/1024x768.jpg")
	if err != nil {
		log.Fatalf("Unable to create testing content, /tmp/1024x768.jpg. %v\n", err)
		os.Exit(-1)
	}

	tmp.Write(getJpeg())
	tmp.Close()
}

type expected struct {
	Size      *PixelDim
	ImgType   string
	MimeType  string
	UseConfig bool
	Link      string
}

func TestSource(t *testing.T) {
	var once sync.Once

	cases := map[string]*expected{
		"http://" + test_server + "/" + file_name: &expected{&PixelDim{1024, 768}, "jpeg", "image/jpeg", false, test_server + "/" + file_name},
		"":                          &expected{Size: nil, UseConfig: false},
		"file://":                   &expected{Size: nil, UseConfig: false}, //file without 'root' option defined
		"localhost://../etc/passwd": &expected{Size: nil, UseConfig: false}, //wrong characters
		"file:// ":                  &expected{Size: nil, UseConfig: true},
		"file://1024x768.jpg":       &expected{&PixelDim{1024, 768}, "jpeg", "image/jpeg", true, "/tmp/1024x768.jpg"},
	}

	for option, want := range cases {
		if want != nil && want.UseConfig {
			once.Do(func() {
				cfg := config.Get()
				cfg.Sources.File.Default = false
				cfg.Sources.File.Root = "/tmp"
			})
		}

		src := Construct(new(Source), option).(*Source)
		if want.Size == nil && src != nil {
			t.Errorf("Expected nil, got %v\n", src)
			continue
		}
		if src != nil {
			if !reflect.DeepEqual(src.Size(), want.Size) {
				t.Errorf("Expected size is %v, got %v\n", want.Size, src.Size())
			}
			if src.Type() != want.ImgType {
				t.Errorf("Expected image type is %v, got %v\n", want.ImgType, src.Type())
			}
			if src.Mime() != want.MimeType {
				t.Errorf("Expected mime type is %v, got %v", src.Mime(), want.MimeType)
			}
			if src.Link() != want.Link {
				t.Errorf("Expected link is %v, got %v\n", want.Link, src.Link())
			}
		}
	}
}

func TestScale(t *testing.T) {
	var scale *Scale
	srcsize := &PixelDim{Width: 1024, Height: 768}

	cases := map[string]*PixelDim{
		"":     nil,
		"800x": &PixelDim{Width: 800, Height: 600},
		"x600": &PixelDim{Width: 800, Height: 600},
		"600":  &PixelDim{600, 450},
		"0.5":  &PixelDim{512, 384},
	}

	for opt, expected := range cases {
		scale = Construct(new(Scale), opt).(*Scale)

		if expected != nil {
			dstsize := scale.Size(srcsize)
			if !reflect.DeepEqual(expected, dstsize) {
				t.Errorf("Expected %v, got %v", expected, dstsize)
			} else {
				log.Println(expected, "=", dstsize)
			}
		} else {
			if scale != nil {
				t.Errorf("Expected nil, got %v\n", scale)
			}
		}
	}
}

func TestCrop(t *testing.T) {
	var crop *Crop

	original := &PixelDim{Width: 1024, Height: 768}

	cases := map[string]*Crop{
		"":               nil,
		"1,1,500,500":    &Crop{param: "user", Roi: &ROI{1, 1, 500, 500}},
		"left,500,500":   &Crop{param: "left", Roi: &ROI{0, 0, 500, 500}},
		"right,500,500":  &Crop{param: "right", Roi: &ROI{524, 0, 500, 500}},
		"bleft,500,500":  &Crop{param: "bleft", Roi: &ROI{0, 268, 500, 500}},
		"bright,500,500": &Crop{param: "bright", Roi: &ROI{524, 268, 500, 500}},
		"center,500,500": &Crop{param: "center", Roi: &ROI{262, 134, 500, 500}},
	}

	for opt, expected := range cases {
		crop = Construct(new(Crop), opt).(*Crop)
		if expected == nil {
			if crop != nil {
				t.Errorf("Expected nil, got %v\n", crop)
			}
		} else {
			if expected.param != crop.param {
				t.Errorf("Expected param is %v, got %v\n", expected.param, crop.param)
			}

			crop.Calc(original)

			if !reflect.DeepEqual(expected.Roi, crop.Roi) {
				t.Errorf("Expected Roi is %v, got %v\n", expected.Roi, crop.Roi)
			}
		}
	}
}
