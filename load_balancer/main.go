package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

var backends = []string{
	"http://backend1:8021", // Backend 1
	"http://backend2:8022", // Backend 2
}

var currentIndex uint64

func lbHandler(w http.ResponseWriter, r *http.Request) {
	index := atomic.AddUint64(&currentIndex, 1) % uint64(len(backends))
	target := backends[index]

	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
}

func main() {
	log.Println("Load Balancer started on :8020")
	http.HandleFunc("/", lbHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8020", nil))
}
