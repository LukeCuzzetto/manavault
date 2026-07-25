package config

import "os"

const defaultPort = "8080"

type Config struct {
	Address string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return Config{
		Address: ":" + port,
	}
}
