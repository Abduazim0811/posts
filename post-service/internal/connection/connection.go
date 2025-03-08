package connection

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"post-service/internal/config"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Connection struct {
	DB    *sql.DB
	Redis *redis.Client
}

func OpenSql() (*sql.DB, error) {
	c := config.Configuration()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBname,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("Failed to open database:", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Println("Unable to connect to database:", err)
		return nil, err
	}

	log.Println("PostgreSQL successfully connected")
	return db, nil
}

func OpenRedis(addr, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Println("Unable to connect to Redis:", err)
		return nil, err
	}

	log.Println("Redis successfully connected")
	return client, nil
}

func NewConnection() (*Connection, error) {
	db, err := OpenSql()
	if err != nil {
		return nil, err
	}
	c := config.Configuration()
	redisClient, err := OpenRedis(c.Redis.Addr, c.Redis.Password) 
	if err != nil {
		db.Close() 
		return nil, err
	}

	return &Connection{
		DB:    db,
		Redis: redisClient,
	}, nil
}

func (c *Connection) Close() {
	if err := c.DB.Close(); err != nil {
		log.Println("Failed to close database:", err)
	}
	if err := c.Redis.Close(); err != nil {
		log.Println("Failed to close Redis:", err)
	}
}