package api

import (
	"hanacaraka/internal/service"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/hello", service.HelloHandler)
}
