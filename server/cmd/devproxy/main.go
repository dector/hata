package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	listenAddr := flag.String("listen", "127.0.0.1:4500", "Address to bind proxy to")
	targetURL := flag.String("target", "http://127.0.0.1:4501", "Upstream app URL")
	flag.Parse()

	target, err := url.Parse(*targetURL)
	if err != nil {
		log.Fatalf("invalid target URL %q: %v", *targetURL, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %v", err)
		http.Error(w, "proxy error", http.StatusBadGateway)
	}

	log.Printf("Local proxy listening on http://%s -> %s", *listenAddr, target.String())
	if err := http.ListenAndServe(*listenAddr, proxy); err != nil {
		log.Fatalf("proxy server failed: %v", err)
	}
}
