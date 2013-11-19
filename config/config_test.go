package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestDefaults(t *testing.T) {
	cfgptr = nil

	if err := ioutil.WriteFile("imagio.conf", []byte("{}"), 0644); err != nil {
		t.Fatalf("Unable to flush imagio.conf. %v\n", err)
	}

	if Get().Method() != METHOD {
		t.Errorf("Expected method is %v, got %v\n", METHOD, Get().Method())
	}

	if Get().Format() != FORMAT {
		t.Errorf("Expected method is %v, got %v\n", FORMAT, Get().Format())
	}

	if Get().CacheSize() != CACHE_SIZE<<20 {
		t.Errorf("Expected cache size is %v, got %v\n", CACHE_SIZE<<20, Get().CacheSize())
	}

	if Get().Quality() != QUALITY {
		t.Errorf("Expected quality is %v, got %v\n", QUALITY, Get().Quality())
	}
}

func TestEmbedJson(t *testing.T) {
	cfgptr = nil
	os.Remove("imagio.conf")

	if Get().Method() != METHOD {
		t.Errorf("Expected method is %v, got %v\n", METHOD, Get().Method())
	}
	if Get().Format() != FORMAT {
		t.Errorf("Expected method is %v, got %v\n", FORMAT, Get().Format())
	}

	if Get().CacheSize() != CACHE_SIZE<<20 {
		t.Errorf("Expected cache size is %v, got %v\n", CACHE_SIZE<<20, Get().CacheSize())
	}

	if Get().Quality() != QUALITY {
		t.Errorf("Expected quality is %v, got %v\n", QUALITY, Get().Quality())
	}

	if Get().CacheSelf() != "http://127.0.0.1:9100" {
		t.Errorf("Expected self option is http://127.0.0.1:9100, got %v\n", Get().CacheSelf())
	}

	peers := []string{"http://127.0.0.1:9100"}
	if !reflect.DeepEqual(peers, Get().CachePeers()) {
		t.Errorf("Expected peers option is %v, got %v\n", peers, Get().CachePeers())
	}
}

func TestFileConf(t *testing.T) {
	cfgptr = nil
	var cfg string = `
    {
        "listen" : "darkstar",

        "defaults" : {
            "format"  : "png",
            "method"  : 4,
            "quality" : 80
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
            "self"  : "http://127.0.0.1:8000",

            "peers" : [
                "http://cache01.local:8000",
                "http://cache02.local:8000",
                "http://cache03.local:8000"
            ],

            "size"  :  "1G"
        }
    }`

	os.Remove("imagio.conf")
	if err := ioutil.WriteFile("imagio.conf", []byte(cfg), 0644); err != nil {
		t.Fatalf("Unable to create testing imagio.conf.", err)
	}

	if Get().Listen() != "darkstar" {
		t.Errorf("Expected listen address is 'darkstar', got %v\n", Get().Listen())
	}

	if Get().Method() != 4 {
		t.Errorf("Expected method is 4, got %v\n", Get().Method())
	}

	if Get().Format() != "png" {
		t.Errorf("Expected method is png, got %v\n", Get().Format())
	}

	if Get().CacheSize() != 1<<30 {
		t.Errorf("Expected cache size is %v, got %v\n", 1<<30, Get().CacheSize())
	}

	if Get().CacheSelf() != "http://127.0.0.1:8000" {
		t.Errorf("Expected self option is 'http://127.0.0.1:8000', got %v\n", Get().CacheSelf())
	}

	peers := []string{"http://cache01.local:8000", "http://cache02.local:8000", "http://cache03.local:8000", "http://127.0.0.1:8000"}
	if !reflect.DeepEqual(peers, Get().CachePeers()) {
		t.Errorf("Expected peers option is %v, got %v\n", peers, Get().CachePeers())
	}

	os.Remove("imagio.conf")
}
