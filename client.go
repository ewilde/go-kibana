package gokong

import (
	"github.com/hashicorp/go-hclog"
	"os"
)

var AppLogger := hclog.New(&hclog.LoggerOptions{
	Name:  "go-kibana",
	Level: hclog.LevelFromString(getEnvironmentValueOrDefault("LOG_LEVEL", "INFO")),
})

func getEnvironmentValueOrDefault(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
