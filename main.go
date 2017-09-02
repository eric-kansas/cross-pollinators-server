package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/eric-kansas/cross-pollinators-server/api"
	"github.com/eric-kansas/cross-pollinators-server/configs"
	"github.com/eric-kansas/cross-pollinators-server/db"
)

var httpServer = &http.Server{
	Addr:         configs.Data.Addr,
	ReadTimeout:  1 * time.Second,
	WriteTimeout: 1 * time.Second,
	IdleTimeout:  1 * time.Second,
}

func init() {
	configs.Initialize()
	err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %+v", err)
	}
}

func main() {
	fmt.Printf("Server started with go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)

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
