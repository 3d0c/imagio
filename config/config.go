package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
)

const (
	CACHE_SIZE = int64(512)
	FORMAT     = "jpeg"
	METHOD     = 3
	QUALITY    = 80
	ALPHA      = 0.5
	LISTEN_ON  = "127.0.0.1:15900"
	CACHE_SELF = "http://127.0.0.1:9100"
)

var defaultCfg string = `
{
    "listen" : "127.0.0.1:15900",

    "defaults" : {
        "format"  : "jpeg",
        "method"  : 3,
        "quality" : 80,
        "alpha"   : 0.5
    },

    "source" : {
        "http" : {
            "root"    : "",
            "default" : true
        },

        "file" : {
            "root"   : "",
            "defaut" : false
        }
    },

    "groupcache" : {
        "self"  : "http://127.0.0.1:9100",
        "peers" : [],
        "size"  : "512M"
    }
}
`

type Source struct {
	Root    string `json:"root"`
	Default bool   `json:"default"`
}

type Config struct {
	ListenOn string `json:"listen"`

	Sources struct {
		Http Source `json:"http"`
		File Source `json:"file"`
	} `json:"source"`

	Defaults struct {
		Format  string  `json:"format"`
		Method  int     `json:"method"`
		Quality int     `json:"quality"`
		Alpha   float64 `json:"blend_alpha"`
	} `json:"defaults"`

	GroupCache struct {
		Self  string   `json:"self"`
		Peers []string `json:"peers"`
		Size  string   `json:"size"`
	} `json:"groupcache"`

	Blend struct {
		With string `json:"with"`
		Mask string `json:"mask"`
		Roi  string `json:roi`
	}
}

var cfgptr *Config

func Get() *Config {
	var data []byte
	if cfgptr != nil {
		return cfgptr
	}

	cfgptr = &Config{}

	data, err := ioutil.ReadFile("imagio.conf")
	if err != nil {
		data = []byte(defaultCfg)
	}

	err = json.Unmarshal(data, cfgptr)
	if err != nil {
		log.Println("Unable to read default config.", err)
	}

	return cfgptr
}

func (this *Config) DumpCfg() error {
	b, err := json.MarshalIndent(this, "", " ")
	if err != nil {
		return errors.New("Unable to Marshal current config." + err.Error())
	}

	err = ioutil.WriteFile("imagio.conf", b, 0644)
	if err != nil {
		return errors.New("Unable to write 'imagio.conf'." + err.Error())
	}

	return nil
}

func (this *Config) Listen() string {
	if this.ListenOn == "" {
		return LISTEN_ON
	}

	return this.ListenOn
}

func (this *Config) Scheme() string {
	if this.Sources.File.Default {
		return "file"
	}

	return "http"
}

func (this *Config) Root(scheme string) (string, error) {
	switch scheme {
	case "file":
		return this.RootFile()
		break

	case "http":
		return this.RootHttp()
		break

	default:
		if this.Sources.File.Default {
			return this.RootFile()
		}
	}

	return this.RootHttp()
}

func (this *Config) RootHttp() (string, error) {
	return this.Sources.Http.Root, nil
}

func (this *Config) RootFile() (string, error) {
	r := []rune(this.Sources.File.Root)
	if len(r) == 0 {
		return "", errors.New("Root for 'file source' sould be defined.")
	}

	if r[len(r)-1] != '/' {
		r = append(r, '/')
	}

	return string(r), nil
}

func (this *Config) CacheSelf() string {
	if this.GroupCache.Self != "" {
		return this.GroupCache.Self
	}

	return CACHE_SELF
}

func (this *Config) CachePeers() []string {
	this.GroupCache.Peers = append(this.GroupCache.Peers, this.CacheSelf())

	return this.GroupCache.Peers
}

func (this *Config) CacheSize() int64 {
	if this.GroupCache.Size == "" {
		return CACHE_SIZE << 20
	}

	r := regexp.MustCompile("([0-9]+)(M|G)")
	result := r.FindStringSubmatch(this.GroupCache.Size)
	if len(result) != 3 {
		log.Printf("Wrong size value '%v', using default.\n", result[2])
		return CACHE_SIZE << 20
	}

	val, err := strconv.Atoi(result[1])
	if err != nil {
		log.Printf("Wrong size value '%v'. %v\n", result[1], err)
		return CACHE_SIZE << 20
	}

	if result[2] == "G" {
		return int64(val) << 30
	}

	if result[2] == "M" {
		return int64(val) << 20
	}

	return CACHE_SIZE << 20
}

func (this *Config) Format() string {
	if this.Defaults.Format == "" {
		return FORMAT
	}

	return this.Defaults.Format
}

func (this *Config) Method() int {
	if this.Defaults.Method == 0 {
		return METHOD
	}

	return this.Defaults.Method
}

func (this *Config) Quality() int {
	if this.Defaults.Quality == 0 {
		return QUALITY
	}

	return this.Defaults.Quality
}

func (this *Config) Alpha() float64 {
	if this.Defaults.Alpha == 0 {
		return ALPHA
	}

	return this.Defaults.Alpha
}

func (this *Config) BlendWith(s string) string {
	if s != "" {
		return s
	}

	return this.Blend.With
}

func (this *Config) BlendMask(s string) string {
	if s != "" {
		return s
	}

	return this.Blend.Mask
}

func (this *Config) BlendRoi(s string) string {
	if s != "" {
		return s
	}

	return this.Blend.Roi
}
