package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eric-kansas/cross-pollinators-server/configs"
)

func init() {
	configs.Initialize()
}

func main() {
	fmt.Println("Hello World")
	setupServer()
}

func setupServer() {
	httpServer := &http.Server{
		Addr:         configs.Data.Addr,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  1 * time.Second,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", Hello)
	httpServer.Handler = mux

	log.Fatal(httpServer.ListenAndServe())
}

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello world!!")
}
