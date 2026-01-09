package getenvfield

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Getenv(key string) string {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println(err)
	}

	val := os.Getenv(key)
	if val == "" {
		log.Println("empty val your key:", key)
	}

	return val
}
