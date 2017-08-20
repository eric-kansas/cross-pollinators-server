package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/eric-kansas/cross-pollinators-server/configs"
)

func init() {
	configs.Initialize()
}

func main() {
	fmt.Printf("Server started with go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
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
	mux.HandleFunc("/healthz", HealthzHandler)
	httpServer.Handler = mux

	log.Fatal(httpServer.ListenAndServe())
}

func HealthzHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Success")
}
