package controllers

import (
	db "bkstream/config"
	"bkstream/structs"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
)

var ctx = context.Background()

type SalesReport struct {
	ID         uint    `json:"id"`
	SalesmanID uint    `json:"salesman_id"`
	Product    string  `json:"product"`
	Quantity   int     `json:"quantity"`
	Revenue    float64 `json:"revenue"`
}

func GetDataCustomerRedis(c *fiber.Ctx) error {
	start := time.Now()
	// if err := db.DB.Limit(1000).Find(&customers).Error; err != nil {
	// 	log.Fatal(err)
	// }

	cacheKey := fmt.Sprintf("customer:test") // Dynamic cache key

	// Try to get the report from Redis
	cachedReport, err := db.RedisClient.Get(ctx, cacheKey).Result()

	exists, err := db.RedisClient.Exists(ctx, cacheKey).Result()

	if errors.Is(err, redis.Nil) || exists == 0 {
		// Cache miss - query the database
		customers := []structs.Customer{}

		if err := db.DB.Limit(10000).Find(&customers).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch customer from database")
		}

		// Serialize the result to JSON
		reportJSON, err := json.Marshal(customers)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to serialize sales report")
		}

		// Store the serialized data in Redis with an expiration time (e.g., 1 minute)
		err = db.RedisClient.Set(ctx, cacheKey, reportJSON, 1*time.Minute).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to cache sales report")
		}

		// Return the result
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "success",
			"data":    customers,
		})
	} else if err != nil {
		// Redis error
		return c.Status(fiber.StatusInternalServerError).SendString("Redis error: " + err.Error())
	}

	// Cache hit - deserialize the result
	var report []structs.Customer
	err = json.Unmarshal([]byte(cachedReport), &report)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to deserialize cached sales report")
	}

	t := time.Now()
	elapsed := t.Sub(start)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    report,
		"elapsed": elapsed,
	})
}
