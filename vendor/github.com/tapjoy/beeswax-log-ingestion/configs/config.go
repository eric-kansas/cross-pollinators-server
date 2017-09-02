package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type beeswaxIngesterConfigs struct {
	Addr            string
	CPMHostname     string
	DISHostname     string
	NewRelicAppName string
	NewRelicEnabled bool
	NewRelicID      string
	Verbose         bool
}

// Data contains loaded configs
var Data beeswaxIngesterConfigs

// Initialize reads specific values from env vars that are needed
// to start up the beeswax log ingester service. Defaults will be used
// if env var is not found
func Initialize() {
	// Load in env vars
	Data = beeswaxIngesterConfigs{
		Addr:            fmt.Sprintf(":%s", getEnv("PORT", "32785")),
		CPMHostname:     getEnv("CPM_HOSTNAME", "http://localhost:80/"),
		DISHostname:     getEnv("DIS_HOSTNAME", "http://localhost:5002/"),
		NewRelicAppName: getEnv("NEW_RELIC_APP_NAME", "dev-beeswax-log-ingestion"),
		NewRelicEnabled: getEnvBool("NEW_RELIC_ENABLED", false),
		NewRelicID:      getEnv("NEW_RELIC_LICENSE_KEY", "abcdefghijklmnopqrstuvwxyzabcdefghijklmn"),
		Verbose:         false,
	}

	log.Printf("BeeswaxIngesterConfigs - %+v\n", Data)
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
