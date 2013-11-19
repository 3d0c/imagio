package query

import (
	"log"
	"strconv"
	"strings"
)

const (
	USER_DEF = 4
	SHORTCUR = 3
)

type ROI struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Crop struct {
	initparam string
	Roi       *ROI
	param     string
	calc      func(x, y, w, h int) *ROI
}

var calcHandlers = map[string]func(int, int, int, int) *ROI{
	"left": func(x, y, w, h int) *ROI {
		return &ROI{0, 0, w, h}
	},

	"right": func(x, y, w, h int) *ROI {
		return &ROI{x - w, 0, w, h}
	},

	"bleft": func(x, y, w, h int) *ROI {
		return &ROI{0, y - h, w, h}
	},

	"bright": func(x, y, w, h int) *ROI {
		return &ROI{x - w, y - h, w, h}
	},

	"center": func(x, y, w, h int) *ROI {
		return &ROI{(x - w) / 2, (y - h) / 2, w, h}
	},
}

func (*Crop) Construct(i ...interface{}) *Crop {
	if len(i) != 1 {
		log.Printf("Wrong arguments count = %d. Expecting 1\n", len(i))
		return nil
	}

	if i[0].([]interface{})[0].(string) == "" {
		return nil
	}

	v := i[0].([]interface{})[0].(string)

	this := &Crop{initparam: v}

	parts := strings.Split(v, ",")

	switch len(parts) {
	case USER_DEF:
		x, err := strconv.Atoi(parts[0])
		y, err := strconv.Atoi(parts[1])
		w, err := strconv.Atoi(parts[2])
		h, err := strconv.Atoi(parts[3])

		this.Roi = &ROI{x, y, w, h}

		if err != nil {
			log.Printf("Illegal crop parameters. x,y,width,height should be integers, `%s` given. Error: %v\n", v, err)
			return nil
		}

		this.param = "user"
		break

	case SHORTCUR:
		this.param = parts[0]

		w, err := strconv.Atoi(parts[1])
		h, err := strconv.Atoi(parts[2])

		if err != nil {
			log.Println("Illegal crop parameters. width and height should be integers, `%s` given. Error: %v\n", v, err)
			return nil
		}

		this.Roi = &ROI{0, 0, w, h}

		var found bool
		if this.calc, found = calcHandlers[this.param]; !found {
			log.Printf("Illegal crop shortcut `%s`. Crop option is `%v`\n", this.param, v)
			return nil
		}

		break

	default:
		return nil
	}

	return this
}

// calculate actual X,Y for crop ROI from original image dimensions
func (this *Crop) Calc(orig *PixelDim) *ROI {
	if this.param == "user" {
		return this.Roi
	}

	this.Roi = this.calc(orig.Width, orig.Height, this.Roi.Width, this.Roi.Height)

	return this.Roi
}
