package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func getRedisConfig(prefix string) (*RedisConfig, error) {
	addr := os.Getenv(prefix + "_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("environment variable %s_ADDR is required", prefix)
	}

	password := os.Getenv(prefix + "_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("environment variable %s_PASSWORD is required", prefix)
	}

	dbStr := os.Getenv(prefix + "_DB")
	if dbStr == "" {
		dbStr = "0" // Default to DB 0 if not specified
	}

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, fmt.Errorf("invalid %s_DB value: %s", prefix, dbStr)
	}

	return &RedisConfig{
		Addr:     addr,
		Password: password,
		DB:       db,
	}, nil
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	ctx := context.Background()

	// Load source Redis configuration
	srcConfig, err := getRedisConfig("SOURCE")
	if err != nil {
		fmt.Printf("Error loading source configuration: %v\n", err)
		os.Exit(1)
	}

	// Load target Redis configuration
	dstConfig, err := getRedisConfig("TARGET")
	if err != nil {
		fmt.Printf("Error loading target configuration: %v\n", err)
		os.Exit(1)
	}

	// Source: Redis
	src := redis.NewClient(&redis.Options{
		Addr:     srcConfig.Addr,
		Password: srcConfig.Password,
		DB:       srcConfig.DB,
	})

	// Target: Redis
	dst := redis.NewClient(&redis.Options{
		Addr:     dstConfig.Addr,
		Password: dstConfig.Password,
		DB:       dstConfig.DB,
	})

	// Test connections
	if err := src.Ping(ctx).Err(); err != nil {
		fmt.Printf("Failed to connect to source Redis: %v\n", err)
		os.Exit(1)
	}

	if err := dst.Ping(ctx).Err(); err != nil {
		fmt.Printf("Failed to connect to target Redis: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Connected to source Redis at %s (DB: %d)\n", srcConfig.Addr, srcConfig.DB)
	fmt.Printf("Connected to target Redis at %s (DB: %d)\n", dstConfig.Addr, dstConfig.DB)

	// Scan all keys
	iter := src.Scan(ctx, 0, "*", 500).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := src.Dump(ctx, key).Result()
		if err != nil {
			fmt.Println("dump error:", err)
			continue
		}
		ttl, _ := src.TTL(ctx, key).Result()

		// Restore into Redis
		err = dst.RestoreReplace(ctx, key, ttl, val).Err()
		if err != nil {
			fmt.Println("restore error:", err)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	fmt.Println("Migration finished 🚀")
}
