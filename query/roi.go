package query

import (
	"log"
	"strconv"
	"strings"
)

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Roi struct {
	InitArea *Rect
	calc     func(x, y, w, h int) *Rect
}

var handlers = map[string]func(int, int, int, int) *Rect{
	"left": func(x, y, w, h int) *Rect {
		return &Rect{0, 0, w, h}
	},

	"right": func(x, y, w, h int) *Rect {
		return &Rect{x - w, 0, w, h}
	},

	"bleft": func(x, y, w, h int) *Rect {
		return &Rect{0, y - h, w, h}
	},

	"bright": func(x, y, w, h int) *Rect {
		return &Rect{x - w, y - h, w, h}
	},

	"center": func(x, y, w, h int) *Rect {
		return &Rect{(x - w) / 2, (y - h) / 2, w, h}
	},
}

// ->x,y,w,h     4
//    this.InitArea{x,y,w,h}
//    this.calc = nil
//    <- InitArea
// ->center,w,h  3
//    this.InitArea{0,0,w,h}
//    this.calc = handlers["center"]
//    <- calculated x,y from source image (w,h are user defined)
// ->x,y         2
//    this.InitArea{x,y,0,0}
//    this.calc = nil
//    <- InitArea (w,h are 0)
//
func (*Roi) Construct(i ...interface{}) *Roi {
	var found bool

	if len(i) != 1 {
		log.Printf("Wrong arguments count = %d. Expecting 1\n", len(i))
		return nil
	}

	if i[0].([]interface{})[0].(string) == "" {
		return nil
	}

	v := i[0].([]interface{})[0].(string)

	this := &Roi{calc: nil}

	parts := strings.Split(v, ",")

	switch len(parts) {
	case 4:
		x, err := strconv.Atoi(parts[0])
		y, err := strconv.Atoi(parts[1])
		w, err := strconv.Atoi(parts[2])
		h, err := strconv.Atoi(parts[3])

		this.InitArea = &Rect{x, y, w, h}

		if err != nil {
			log.Printf("Illegal parameters: x,y,width,height should be integers, `%s` given. Error: %v\n", v, err)
			return nil
		}

		break

	case 3:
		w, err := strconv.Atoi(parts[1])
		h, err := strconv.Atoi(parts[2])

		if err != nil {
			log.Println("Illegal parameters. width and height should be integers, `%s` given. Error: %v\n", v, err)
			return nil
		}

		this.InitArea = &Rect{0, 0, w, h}

		if this.calc, found = handlers[parts[0]]; !found {
			log.Printf("Illegal roi shortcut `%s`\n", parts[0])
			return nil
		}

		break

	case 2:
		x, err := strconv.Atoi(parts[0])
		y, err := strconv.Atoi(parts[1])

		if err != nil {
			log.Printf("Illegal parameters. x,y should be integers, `%s` given. Error: %v\n", v, err)
			return nil
		}

		this.InitArea = &Rect{x, y, 0, 0}

		break

	default:
		return nil
	}

	return this
}

func (this *Roi) Calc(orig *PixelDim) *Rect {
	if this.calc == nil {
		return this.InitArea
	}

	return this.calc(orig.Width, orig.Height, this.InitArea.Width, this.InitArea.Height)
}
