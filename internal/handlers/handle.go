package handlers

import (
	"encoding/json"
	"net/http"

	"676f.dev/rfc7807"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf"
)

type Handler struct {
	rc *redis.Client
	k  *koanf.Koanf
}

func NewHandler(rc *redis.Client, k *koanf.Koanf) *Handler {
	return &Handler{
		rc: rc,
		k:  k,
	}
}

func errorResponse(c *fiber.Ctx, status int, title string, detail string) {
	problem := rfc7807.NewError(title).
		SetStatus(status).
		SetDetail(detail)

	responseBody, err := json.Marshal(problem)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(status)
	_, err = c.Write(responseBody)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
}
