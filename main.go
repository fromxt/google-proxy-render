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

	// 
	proxyPassword := os.Getenv("PROXY_PASSWORD")

	// 
	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 
		if proxyPassword == "" {
			p.ServeHTTP(w, r)
			return
		}

		// 
		clientPwd := r.URL.Query().Get("pwd")
		
		// 
		if clientPwd == "" {
			clientPwd = r.Header.Get("X-Proxy-Password")
		}

		// 
		if subtle.ConstantTimeCompare([]byte(clientPwd), []byte(proxyPassword)) != 1 {
			http.Error(w, "Unauthorized: Incorrect or missing password", http.StatusUnauthorized)
			return
		}

		//
		p.ServeHTTP(w, r)
	})

	// 
	http.ListenAndServe(":"+os.Getenv("PORT"), authHandler)
}
