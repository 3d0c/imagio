package imgproc

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	. "github.com/3d0c/imagio/query"
	. "github.com/3d0c/imagio/utils"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"
)

const test_server = "localhost:41899"
const file_name = "1024x768.jpg"
const sample_jpeg = "1f8b080866f18152020331303234783736382e6a706700edce3b0ec2301045d13718210a8a588226280d8515576e421b2b418a050bcd3a2858049f869d1807d1f069d2a277dc5dc93313cff18eecd0ed3b8800921ee20d3bacf3dc16b631a6699d736df075edc308c31029acad4ce5cbd207bf1df5fd35e4083d579842c906132d4a4b3c61f53cf5cd2ce565f6593154f959f577bd62a1246d511a1e3d8888888888888888e89f49bc3c005c75237d1b130000"

func serveJpeg() {
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, file_name, time.Now(), bytes.NewReader(data))
	})

	log.Fatal(http.ListenAndServe(test_server, nil))
}

func init() {
	go serveJpeg()
}

type expected struct {
	Size    *PixelDim
	ImgType string
}

func TestProc(t *testing.T) {
	cases := map[*Options]*expected{
		&Options{
			Format:  "jpg",
			Quality: 80,
			Method:  3,
			Base:    Construct(new(Source), "http://"+test_server+"/"+file_name).(*Source),
			Scale:   Construct(new(Scale), "100x").(*Scale),
		}: &expected{Size: &PixelDim{100, 75}, ImgType: "jpeg"},
		&Options{
			Format:  "jpg",
			Quality: 80,
			Base:    Construct(new(Source), "http://"+test_server+"/"+file_name).(*Source),
			CropRoi: Construct(new(Roi), "1,1,500,500").(*Roi),
			Scale:   nil,
		}: &expected{&PixelDim{500, 500}, "jpeg"},
		&Options{
			Format:  "jpg",
			Quality: 80,
			Method:  3,
			Base:    Construct(new(Source), "http://"+test_server+"/"+file_name).(*Source),
			Scale:   Construct(new(Scale), "100x").(*Scale),
			CropRoi: Construct(new(Roi), "center,500,500").(*Roi),
		}: &expected{&PixelDim{100, 100}, "jpeg"},
		&Options{
			Format: "png",
			Method: 3,
			Base:   Construct(new(Source), "http://"+test_server+"/"+file_name).(*Source),
			Scale:  Construct(new(Scale), "100x").(*Scale),
		}: &expected{&PixelDim{100, 75}, "png"},
		&Options{
			Format: "png",
			Method: 3,
			Base:   Construct(new(Source), "http://"+test_server+"/"+file_name).(*Source),
		}: &expected{&PixelDim{1024, 768}, "png"},
	}

	for option, want := range cases {
		b := Do(option)

		if b == nil {
			t.Errorf("Expected data, result is nil\n")
		}

		cfg, imgType, err := image.DecodeConfig(bytes.NewReader(b))
		if err != nil {
			t.Error(err)
		}

		resultSize := &PixelDim{cfg.Width, cfg.Height}

		if !reflect.DeepEqual(resultSize, want.Size) {
			t.Errorf("Expected size is %v, got %v\n", want.Size, resultSize)
		}

		if imgType != want.ImgType {
			t.Errorf("Expected image type is %v, got %v", want.ImgType, imgType)
		}
	}
}
