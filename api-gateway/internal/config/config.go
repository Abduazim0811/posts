package config

import "os"

type Config struct {
	ApiGateway struct {
		Port string
	}
	User struct {
		Port string
	}
	Post struct {
		Port string
	}
}

func Configuration() *Config {
	c := &Config{}
	c.ApiGateway.Port = osGetenv("API_GATEWAY", "api_gateway:7777")

	c.User.Port = osGetenv("USER_PORT", "user_service:8888")

	c.Post.Port = osGetenv("POST_PORT", "post_service:9999")
	return c
}

func osGetenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
