package config

import "os"

type Config struct {
	Database struct {
		User     string
		Password string
		Host     string
		Port     string
		DBname   string
	}
	Post  struct {
		Host string
		Port string
	}
	Redis struct {
		Addr     string
		Password string
	}
}

func Configuration() *Config {
	c := &Config{}

	c.Database.User = osGetenv("DB_USER", "postgres")
	c.Database.Password = osGetenv("DB_PASSWORD", "Abdu0811")
	c.Database.Host = osGetenv("DB_HOST", "localhost")
	c.Database.Port = osGetenv("DB_PORT", "5432")
	c.Database.DBname = osGetenv("DB_NAME", "postservice")

	c.Post.Host = osGetenv("POST_HOST", "tcp")
	c.Post.Port = osGetenv("POST_PORT", "localhost:9999")

	c.Redis.Addr = osGetenv("REDIS_ADDR", "localhost:6379")
	c.Redis.Password = osGetenv("REDIS_PASSWORD", "")
	return c
}

func osGetenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}