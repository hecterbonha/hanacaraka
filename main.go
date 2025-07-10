package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tm := time.Now().Format(time.RFC1123)
	response := map[string]string{"time": tm, "data": r.Pattern}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func main() {
	mux := http.NewServeMux()

	th := http.HandlerFunc(timeHandler)

	mux.Handle("/time", th)

	log.Print("Listening...")
	http.ListenAndServe(":3000", mux)
}
