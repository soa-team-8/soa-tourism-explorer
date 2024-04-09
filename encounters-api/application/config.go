package application

import (
	"os"
	"strconv"
)

type Config struct {
	PostgresAddress string
	MongoDBAddress  string
	ServerPort      uint16
}

func LoadConfig() Config {
	cfg := Config{
		PostgresAddress: "postgres://postgres:super@localhost:5432/encounters",
		MongoDBAddress:  "mongodb://localhost:27017",
		ServerPort:      3030,
	}

	if postgresAddr, exists := os.LookupEnv("POSTGRES_ADDR"); exists {
		cfg.PostgresAddress = postgresAddr
	}

	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	return cfg
}
