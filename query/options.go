package query

import (
	"github.com/3d0c/imagio/config"
	. "github.com/3d0c/imagio/utils"
	"log"
	"net/url"
	"reflect"
	"strconv"
)

var supportedOptions = map[string]interface{}{
	"jpeg": "jpeg", "jpg": "jpeg", "png": "png", "gif": "gif", "json": "json",
	"NN": 1, "LINEAR": 2, "CUBIC": 3, "AREA": 4, "LANCZOS": 5,
	"true": true, "false": false, "alpha": 0.5,
}

type Options struct {
	Base       *Source
	Scale      *Scale
	CropRoi    *Roi
	Format     string
	Method     int
	Quality    int
	Alpha      float64
	Foreground *Source
	Mask       *Source
	BlendRoi   *Roi
}

func (*Options) Construct(i ...interface{}) *Options {
	if len(i) != 1 {
		log.Printf("Wrong arguments count = %d. Expecting 1\n", len(i))
		return nil
	}

	v := i[0].([]interface{})[0]

	switch v.(type) {
	case *url.URL:
		return parseQuery(v.(*url.URL))

	case string:
		u, err := url.Parse(v.(string))
		if err != nil {
			log.Println("Unable to parse query:", v)
			return nil
		}

		return parseQuery(u)
		break

	default:
		log.Println("Unsupported type:", reflect.TypeOf(v))
	}

	return nil
}

func parseQuery(u *url.URL) *Options {
	query := u.Query()
	log.Println("in:", u.String())
	this := &Options{
		CropRoi: Construct(new(Roi), query.Get("crop")).(*Roi),
		Scale:   Construct(new(Scale), query.Get("scale")).(*Scale),
		Base:    Construct(new(Source), query.Get("source")).(*Source),

		Format:  get(query.Get("format"), config.Get().Format()).(string),
		Method:  get(query.Get("method"), config.Get().Method()).(int),
		Alpha:   getFloat(query.Get("blend_alpha"), config.Get().Alpha()),
		Quality: getInt(query.Get("quality"), config.Get().Quality()),

		Foreground: Construct(new(Source), config.Get().BlendWith(query.Get("blend_with"))).(*Source),
		Mask:       Construct(new(Source), config.Get().BlendMask(query.Get("blend_mask"))).(*Source),
		BlendRoi:   Construct(new(Roi), config.Get().BlendRoi(query.Get("blend_roi"))).(*Roi),
	}

	return this
}

func get(key string, def interface{}) interface{} {
	if val, found := supportedOptions[key]; found {
		return val
	}

	return def
}

func getFloat(key string, def float64) float64 {
	if val, err := strconv.ParseFloat(key, 64); err == nil {
		return val
	}

	return def
}

func getInt(key string, def int) int {
	if val, err := strconv.Atoi(key); err == nil {
		return val
	}

	return def
}
