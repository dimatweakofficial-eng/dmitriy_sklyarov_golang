package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Dsn   DsnConfig
	Token TokenConfig
}

type DsnConfig struct {
	DsnName string
}

type TokenConfig struct {
	TokenName string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Env file not found, using system environment variables")
	}

	return &Config{
		Dsn: DsnConfig{
			DsnName: os.Getenv("DSN"),
		},
		Token: TokenConfig{
			TokenName: os.Getenv("SECRET"),
		},
	}
}
