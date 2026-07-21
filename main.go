package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	// 1. Set the target server URL (DuckDuckGo official homepage)
	target := "https://duckduckgo.com"
	remote, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	// 2. Create the reverse proxy instance
	proxy := httputil.NewSingleHostReverseProxy(remote)

	// 3. Modify request headers: Rewrite the Host header to match the target domain
	// to prevent the target server from rejecting the request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = remote.Host
	}

	// 4. Modify response headers: Remove security headers (CSP / X-Frame-Options)
	// that cause blank screens in the browser
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Content-Security-Policy")
		resp.Header.Del("Content-Security-Policy-Report-Only")
		resp.Header.Del("X-Frame-Options")
		resp.Header.Del("Strict-Transport-Security")
		return nil
	}

	// 5. Handle normal interruption logs: Ignore errors caused by client-side cancellations/refreshes
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if errors.Is(err, context.Canceled) {
			return // Ignore client-initiated cancellations quietly without logging
		}
		log.Printf("[Proxy Error] %v", err)
		w.WriteHeader(http.StatusBadGateway)
	}

	// 6. Read the PORT environment variable injected by Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port for local testing
	}

	log.Printf("Reverse proxy started, listening on port :%s, target: %s\n", port, target)

	// 7. Start the HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
