package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	config "github.com/GmbLucas/front-entry/pkg"
)

var exit = make(chan error)

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(res, req)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if t, ok := config.Get().Proxy[r.Host]; ok {
		serveReverseProxy(t, w, r)
		log.Println("OK", r.RemoteAddr, r.Host, r.URL)
	} else {
		log.Println("KO", r.RemoteAddr, r.Host, r.URL)
	}
}

func main() {
	if err := config.Init(); err != nil {
		fmt.Print("config init: ", err)
		return
	}

	http.HandleFunc("/", handler)
	go func() {
		if err := http.ListenAndServe(":80", nil); err != nil {
			exit<-err
		}
	}()
	go func() {
		if err := http.ListenAndServeTLS(":443", config.Get().Tls["certfile"], config.Get().Tls["keyfile"], nil); err != nil {
			exit<-err
		}
	}()

	log.Println("ready")
	log.Fatal(<-exit)
}