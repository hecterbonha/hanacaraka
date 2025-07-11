package service

import (
	"fmt"
	"net/http"
)


func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Hello from native Go backend!")
}
