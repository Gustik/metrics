package handler

import (
	"net/http"
)

func SetupRouter(handler *MetricHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update/counter/{metric}/{value}", handler.UpdateCounter)
	mux.HandleFunc("GET /value/counter/{metric}", handler.GetCounter)

	mux.HandleFunc("POST /update/gauge/{metric}/{value}", handler.UpdateGauge)
	mux.HandleFunc("GET /value/gauge/{metric}", handler.GetGauge)

	// Для всех остальных типов
	mux.HandleFunc("POST /update/{type}/{metric}/{value}", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
	})

	mux.HandleFunc("GET /value/{type}/{metric}", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
	})

	return ContentTypeMiddleware("text/plain; charset=utf-8")(mux)
}

func ContentTypeMiddleware(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			next.ServeHTTP(w, r)
		})
	}
}
