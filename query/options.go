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
	"true": true, "false": false,
}

type Options struct {
	Source  *Source
	Scale   *Scale
	Crop    *Crop
	Format  string
	Method  int
	Quality int
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

	this := &Options{
		Crop:    Construct(new(Crop), query.Get("crop")).(*Crop),
		Scale:   Construct(new(Scale), query.Get("scale")).(*Scale),
		Source:  Construct(new(Source), query.Get("source")).(*Source),
		Format:  get(query.Get("format"), config.Get().Format()).(string),
		Method:  get(query.Get("method"), config.Get().Method()).(int),
		Quality: getInt(query.Get("quality"), config.Get().Quality()),
	}

	return this
}

func get(key string, def interface{}) interface{} {
	if val, found := supportedOptions[key]; found {
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
