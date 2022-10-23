package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/michaelpeterswa/ratelimit-demo/internal/handlers"
	"github.com/michaelpeterswa/ratelimit-demo/internal/kv"
	"go.uber.org/zap"
)

func main() {
	k := koanf.New(".")

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	if err := k.Load(file.Provider(os.Getenv("CONFIG_FILE")), yaml.Parser()); err != nil {
		logger.Fatal("error loading config file", zap.Error(err))
	}

	redisClient := kv.NewRedisClient(k.String("redis.url"), k.String("redis.port"))
	h := handlers.NewHandler(redisClient.Client, k)

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers:
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("Strict-Transport-Security", "max-age=5184000")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-DNS-Prefetch-Control", "off")

		// Go to next middleware:
		return c.Next()
	})

	_ = app.Use(h.RatelimitHandler)
	_ = app.Post("/", h.IDHandler)

	err = app.Listen(":8080")
	if err != nil {
		logger.Fatal("error starting server", zap.Error(err))
	}
}
