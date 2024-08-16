package common

import (
	"os"
)

func EnvConfig(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
