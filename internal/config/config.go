// config.go
package config

import "os"

type Config struct {
	MongoURI  string
	AuthURL   string
	OrdersURL string
	Port      string
}

func Load() *Config {
	return &Config{
		MongoURI:  getEnv("MONGO_URI", "mongodb://host.docker.internal:27017"),
		AuthURL:   getEnv("AUTH_URL", "http://host.docker.internal:3000"),
		OrdersURL: getEnv("ORDERS_URL", "http://host.docker.internal:3004"),
		Port:      getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
