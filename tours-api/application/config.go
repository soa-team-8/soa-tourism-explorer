package application

import (
	"os"
	"strconv"
)

type Config struct {
	PostgresAddress string
	ServerPort      uint16
}

func LoadConfig() Config {
	config := Config{
		PostgresAddress: "postgres://postgres:super@tours_database:5432/tours",
		ServerPort:      3000,
	}

	if postgresAddr, exists := os.LookupEnv("POSTGRES_ADDR"); exists {
		config.PostgresAddress = postgresAddr
	}

	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			config.ServerPort = uint16(port)
		}
	}

	return config
}
