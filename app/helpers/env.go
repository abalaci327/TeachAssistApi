package helpers

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnvironment() {
	if os.Getenv("APP_ENV") == "DEV" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}
}
