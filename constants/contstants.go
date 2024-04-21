package constants

import "os"

var (
	PORT        = "8080"
	JWT_SECRET  = "SECRET123"
	CONTEXT_KEY = "context"
)

func init() {
	if envPort, exists := os.LookupEnv("PORT"); exists {
		PORT = envPort
	}
	if envJwtSecret, exists := os.LookupEnv("JWT_SECRET"); exists {
		JWT_SECRET = envJwtSecret
	}
	if envContextKey, exists := os.LookupEnv("CONTEXT_KEY"); exists {
		CONTEXT_KEY = envContextKey
	}
}
