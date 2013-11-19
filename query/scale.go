package query

import (
	"log"
	"math"
	"strconv"
	"strings"
)

const (
	MAX_RATIO = 6
	MIN_RATIO = 0.08
)

type PixelDim struct {
	Width  int
	Height int
}

type Scale struct {
	width  int
	height int
	maxdim float64
}

func (*Scale) Construct(i ...interface{}) *Scale {
	if len(i) != 1 {
		log.Printf("Wrong arguments count = %d. Expecting 1.\n", len(i))
		return nil
	}

	v := i[0].([]interface{})[0].(string)
	if v == "" {
		log.Printf("Illegal scale option '%v'\n", v)
		return nil
	}

	parts := strings.Split(v, "x")

	switch len(parts) {
	case 1:
		f, err := strconv.ParseFloat(parts[0], 32)
		if err != nil {
			log.Printf("Illegal scale option '%v'\n", v)
			return nil
		}

		return &Scale{maxdim: f}

	case 2:
		w, _ := strconv.Atoi(parts[0])
		h, _ := strconv.Atoi(parts[1])

		return &Scale{
			width:  w,
			height: h,
		}

	default:
		log.Println("Illegal scale option '%v'\n", v)
		return nil
	}

	log.Printf("Illegal size value = '%v'\n", v)

	return nil
}

func (this *Scale) Size(src *PixelDim) *PixelDim {
	ratio := float64(src.Width) / float64(src.Height)

	if ratio > MAX_RATIO || ratio < MIN_RATIO {
		log.Printf("Aspect ratio is %.4f. Seems unreasonable to scale.\n", ratio)
		return nil
	}

	if this.height > 0 && this.width == 0 {
		return &PixelDim{
			Width:  int(math.Floor(float64(this.height) * ratio)),
			Height: this.height,
		}
	}

	if this.height == 0 && this.width > 0 {
		return &PixelDim{
			Width:  this.width,
			Height: int(math.Floor(float64(this.width) / ratio)),
		}
	}

	if this.maxdim > 0 && this.maxdim < 1 {
		return &PixelDim{
			Width:  int(math.Floor(float64(src.Width) * this.maxdim)),
			Height: int(math.Floor(float64(src.Height) * this.maxdim)),
		}
	}

	if this.maxdim >= 1 && src.Width > src.Height {
		this.width = int(this.maxdim)
		return this.Size(src)
	}

	if this.maxdim >= 1 && src.Width < src.Height {
		this.height = int(this.maxdim)
		return this.Size(src)
	}

	return &PixelDim{Width: src.Width, Height: src.Height}
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}
