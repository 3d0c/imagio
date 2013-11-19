package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/3d0c/imagio/config"
	"github.com/3d0c/imagio/imgproc"
	"github.com/3d0c/imagio/query"
	. "github.com/3d0c/imagio/utils"
	"github.com/golang/groupcache"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var cacheGroup *groupcache.Group

func initCacheGroup() {
	self := config.Get().CacheSelf()

	pool := groupcache.NewHTTPPool(self)
	pool.Set(config.Get().CachePeers()...)

	if self != "" {
		log.Println("Cache listen on:", strings.TrimLeft(self, "http://"))
		go http.ListenAndServe(strings.TrimLeft(self, "http://"), http.HandlerFunc(pool.ServeHTTP))
	}

	cacheGroup = groupcache.NewGroup("imagio-storage", config.Get().CacheSize(), groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			dest.SetBytes(imgproc.Do(
				Construct(new(query.Options), key).(*query.Options),
			))
			return nil
		}),
	)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	dumpcfg := flag.Bool("dumpcfg", false, "Dump config")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *dumpcfg {
		config.Get().DumpCfg()
		os.Exit(0)
	}

	initCacheGroup()

	log.Printf("Service listen on %v\n", config.Get().Listen())

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			var data []byte
			var ctx groupcache.Context

			cacheGroup.Get(ctx, r.URL.String(), groupcache.AllocatingByteSliceSink(&data))

			http.ServeContent(w, r, r.URL.String(), time.Now(), bytes.NewReader(data))
		},
	)

	http.HandleFunc("/nocache",
		func(w http.ResponseWriter, r *http.Request) {
			var result []byte
			result = imgproc.Do(Construct(new(query.Options), r.URL).(*query.Options))
			w.Write(result)
		},
	)

	http.HandleFunc("/stat",
		func(w http.ResponseWriter, r *http.Request) {
			// awesome stat. not implemented yet.
		},
	)

	log.Fatal(http.ListenAndServe(config.Get().Listen(), nil))
}
