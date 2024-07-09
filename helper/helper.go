package helper

import "github.com/joho/godotenv"

// LoadEnv loads environment variables from the .env file.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("failed to connect to env")
	}
}
