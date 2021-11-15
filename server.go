package main

import (
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/NYTimes/gcs-helper/v3/handlers"
	"github.com/rs/cors"
)

func getHandler(c handlers.Config, client *storage.Client, hc *http.Client) http.HandlerFunc {
	proxyHandler := cors.Default().Handler(handlers.Proxy(c, hc))
	mapHandler := cors.Default().Handler(handlers.Map(c, client))

	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, c.Proxy.Endpoint):
			r.URL.Path = strings.Replace(r.URL.Path, c.Proxy.Endpoint, "", 1)
			if !strings.HasPrefix(r.URL.Path, "/") {
				r.URL.Path = "/" + r.URL.Path
			}
			proxyHandler.ServeHTTP(w, r)
		case strings.HasPrefix(r.URL.Path, c.Map.Endpoint):
			r.URL.Path = strings.Replace(r.URL.Path, c.Map.Endpoint, "", 1)
			mapHandler.ServeHTTP(w, r)
		case r.URL.Path == "/":
			// healthcheck
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}
}
