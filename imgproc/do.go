package imgproc

//#cgo pkg-config: opencv
//#include "cv_handler.h"
import "C"

import (
	. "github.com/3d0c/imagio/query"
	"unsafe"
)

type Blob _Ctype_Blob

func Do(o *Options) []byte {
	if o.Source == nil {
		return nil
	}

	if o.Crop != nil && o.Scale != nil {
		// if both options selected, crop will be first, the scale size will be calculated from cropped dimension
		roi := o.Crop.Calc(o.Source.Size())
		zoom := o.Scale.Size(&PixelDim{roi.Width, roi.Height})

		return cvHandler(o, zoom, roi)
	}

	if o.Crop != nil {
		return cvHandler(o, nil, o.Crop.Calc(o.Source.Size()))
	}

	if o.Scale != nil {
		return cvHandler(o, o.Scale.Size(o.Source.Size()), nil)
	}

	return cvHandler(o, o.Source.Size(), nil)
}

func cvHandler(o *Options, zoom *PixelDim, roi *ROI) []byte {
	var data []byte

	blob := &Blob{
		data:   (*C.uchar)(unsafe.Pointer(&o.Source.Blob()[0])),
		length: C.uint(len(o.Source.Blob())),
	}

	result := C.cv_handler(
		(*_Ctype_Blob)(unsafe.Pointer(blob)),
		(*_Ctype_PixelDim)(unsafe.Pointer(zoom)),
		C.int(o.Quality), C.int(o.Method), C.CString("."+o.Format),
		(*_Ctype_ROI)(unsafe.Pointer(roi)),
	)

	if result != nil {
		length := result.length
		data = C.GoBytes(unsafe.Pointer(result.data), C.int(length))
	}

	C.free(unsafe.Pointer(result.data))
	C.free(unsafe.Pointer(result))

	return (data)
}
