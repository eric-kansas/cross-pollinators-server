package main

import (
	"fmt"
	"log"

	"github.com/tapjoy/beeswax-log-ingestion/configs"
	"github.com/tapjoy/beeswax-log-ingestion/ingestion"
)

// TODO: Add config check for production dependencies (ie. new relic)

func init() {
	configs.Initialize()
}

func main() {
	fmt.Println("Running Beeswax Log Ingestion Server!!!")

	opts, err := getServerOptions()
	if err != nil {
		log.Fatalf("Error getting server options: %s", err)
	}

	server, err := ingestion.NewServer(configs.Data.Addr, configs.Data.CPMHostname, opts...)
	if err != nil {
		log.Fatalf("Failed to create ingestion server: %s", err)
	}

	log.Fatal(server.ListenAndServe())
}
