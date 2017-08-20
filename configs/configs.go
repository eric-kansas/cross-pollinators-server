package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type crossPollinatorConfigs struct {
	Addr    string
	Verbose bool
}

// Data contains loaded configs
var Data crossPollinatorConfigs

// Initialize reads specific values from env vars that are needed
// to start up the cross pollinators service. Defaults will be used
// if env var is not found
func Initialize() {
	// Load in env vars
	Data = crossPollinatorConfigs{
		Addr:    fmt.Sprintf(":%s", getEnv("PORT", "3030")),
		Verbose: false,
	}

	log.Printf("Cross Pollinators Service Configs - %+v\n", Data)
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	b, err := strconv.ParseBool(value)
	if err == nil {
		return b
	}

	return defaultValue
}
