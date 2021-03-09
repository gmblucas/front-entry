package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
)

var exit = make(chan error)
var conf Config
type Config struct {
	Proxy map[string]string `toml:proxy`
	Tls map[string]string `toml:tls`
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(res, req)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if t, ok := conf.Proxy[r.Host]; ok {
		serveReverseProxy(t, w, r)
		log.Println("OK", r.RemoteAddr, r.Host, r.URL)
	} else {
		log.Println("KO", r.RemoteAddr, r.Host, r.URL)
	}
}

func main() {
	if _, err := toml.DecodeFile(os.Args[1], &conf); err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	go func() {
		if err := http.ListenAndServe(":80", nil); err != nil {
			exit<-err
		}
	}()
	go func() {
		if err := http.ListenAndServeTLS(
			":443",
			conf.Tls["certfile"],
			conf.Tls["keyfile"],
			nil); err != nil {
			exit<-err
		}
	}()

	panic(<-exit)
}