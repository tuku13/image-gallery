package constants

import "os"

var (
	PORT = "8080"
)

func init() {
	if envPort, exists := os.LookupEnv("PORT"); exists {
		PORT = envPort
	}
}
