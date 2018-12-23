package imgproc

//#cgo pkg-config: opencv
//#cgo CFLAGS: -Wno-error=unused-function
//#cgo LDFLAGS: -lopencv_imgproc -lopencv_core -lopencv_highgui
//#include "cv_handler.h"
import "C"

import (
	. "github.com/3d0c/imagio/query"
	. "github.com/3d0c/imagio/utils"
	"log"
	"unsafe"
)

type Blob _Ctype_Blob
type CvRect _Ctype_CvRect

func Do(o *Options) []byte {
	return Filters(PrimaryActions(o))
}

func PrimaryActions(o *Options) (*Options, []byte) {
	if o.Base == nil {
		return o, nil
	}

	if o.CropRoi != nil && o.Scale != nil {
		// if both options selected, crop will be first, the scale size will be calculated from cropped dimension
		roi := o.CropRoi.Calc(o.Base.Size())
		zoom := o.Scale.Size(&PixelDim{roi.Width, roi.Height})

		return o, resize(o, zoom, roi)
	}

	if o.CropRoi != nil {
		return o, resize(o, nil, o.CropRoi.Calc(o.Base.Size()))
	}

	if o.Scale != nil {
		return o, resize(o, o.Scale.Size(o.Base.Size()), nil)
	}

	return o, resize(o, o.Base.Size(), nil)
}

func Filters(o *Options, b []byte) []byte {
	if b == nil {
		return nil
	}

	if o.Foreground != nil {
		var roi *Rect = nil
		if o.BlendRoi != nil {
			roi = o.BlendRoi.Calc(o.Base.Size())
		}

		return blend(Construct(new(Source), b).(*Source), o, roi)
	}

	return b
}

func resize(o *Options, zoom *PixelDim, roi *Rect) []byte {
	var data []byte

	result := C.resizer(
		(*C.Blob)(unsafe.Pointer(blobptr(o.Base))),
		(*C.PixelDim)(unsafe.Pointer(zoom)),
		C.int(o.Quality), C.int(o.Method), C.CString("."+o.Format),
		(*C.CvRect)(initCvRect(roi)),
	)

	if result != nil {
		length := result.length
		data = C.GoBytes(unsafe.Pointer(result.data), C.int(length))

		C.free(unsafe.Pointer(result.data))
		C.free(unsafe.Pointer(result))
	}

	return (data)
}

func blend(base *Source, o *Options, roi *Rect) []byte {
	var data []byte
	rect := &CvRect{0, 0, 0, 0}

	if roi != nil {
		if w := (roi.X + o.Foreground.Size().Width); w > base.Size().Width {
			log.Printf("Wrong blend_roi: width %d > %d. Using (x = 0).\n", w, base.Size().Width)
			roi.X = 0
		}

		if h := (roi.Y + o.Foreground.Size().Height); h > base.Size().Height {
			log.Printf("Wrong blend_roi: height %d > %d. Using (y = 0).\n", h, base.Size().Height)
			roi.Y = 0
		}

		rect = &CvRect{C.int(roi.X), C.int(roi.Y), C.int(roi.Width), C.int(roi.Height)}
	}

	result := C.blender(
		(*C.Blob)(blobptr(base)),
		(*C.Blob)(blobptr(o.Foreground)),
		(*C.Blob)(blobptr(o.Mask)),
		C.int(o.Quality), C.CString("."+o.Format), C.float(o.Alpha),
		(*C.CvRect)(rect),
	)

	if result != nil {
		length := result.length
		data = C.GoBytes(unsafe.Pointer(result.data), C.int(length))

		C.free(unsafe.Pointer(result.data))
		C.free(unsafe.Pointer(result))
	}

	return (data)
}

func blobptr(s *Source) *Blob {
	if s == nil {
		return nil
	}

	return &Blob{
		data:   (*C.uchar)(unsafe.Pointer(&s.Blob()[0])),
		length: C.uint(len(s.Blob())),
	}
}

func initCvRect(roi *Rect) *CvRect {
	if roi == nil {
		return nil
	}

	if roi.Width == 0 || roi.Height == 0 {
		log.Println("Wrong roi init for crop action, should contain width and height")
		return nil
	}

	return &CvRect{C.int(roi.X), C.int(roi.Y), C.int(roi.Width), C.int(roi.Height)}
}
