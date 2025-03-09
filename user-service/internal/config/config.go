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
	User  struct {
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
	c.Database.Host = osGetenv("DB_HOST", "user_postgres")
	c.Database.Port = osGetenv("DB_PORT", "5432")
	c.Database.DBname = osGetenv("DB_NAME", "userservice")

	c.User.Host = osGetenv("USER_HOST", "tcp")
	c.User.Port = osGetenv("USER_PORT", "user_service:8888")

	c.Redis.Addr = osGetenv("REDIS_ADDR", "redis:6379")
	c.Redis.Password = osGetenv("REDIS_PASSWORD", "")
	return c
}

func osGetenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}