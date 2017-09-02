package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/eric-kansas/cross-pollinators-server/server/api"
	"github.com/eric-kansas/cross-pollinators-server/server/configs"
)

var httpServer = &http.Server{
	Addr:         configs.Data.Addr,
	ReadTimeout:  1 * time.Second,
	WriteTimeout: 1 * time.Second,
	IdleTimeout:  1 * time.Second,
}

func init() {
	fmt.Printf("Server started: Version %s", "2.0.0 \n")
	fmt.Printf("Server running go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	configs.Initialize()
}

func main() {
	setupAPI()

	log.Fatal(httpServer.ListenAndServe())
}

func setupAPI() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", api.HealthzHandler)
	mux.HandleFunc("/login", api.LoginHandler)
	mux.HandleFunc("/register", api.RegisterHandler)
	mux.HandleFunc("/dostuff", api.DoTheThingsHandler)
	httpServer.Handler = mux
}
