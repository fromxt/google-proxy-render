package main

import (
	"crypto/subtle"
	"github.com/AaronO/gogo-proxy"
	"net/http"
	"os"
)

func main() {
	p, _ := proxy.New(proxy.ProxyOptions{
		Balancer: func(req *http.Request) (string, error) {
			return "https://as.priv.au", nil
		},
	})

	proxyPassword := os.Getenv("PROXY_PASSWORD")

	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 
		if proxyPassword == "" {
			p.ServeHTTP(w, r)
			return
		}

		// 
		_, pass, ok := r.BasicAuth()

		// 
		if !ok || subtle.ConstantTimeCompare([]byte(pass), []byte(proxyPassword)) != 1 {
			// 
			w.Header().Set("WWW-Authenticate", `Basic realm="Google Proxy Restricted Access"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		
		p.ServeHTTP(w, r)
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), authHandler)
}
