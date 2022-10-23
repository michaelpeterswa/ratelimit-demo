package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RatelimitHandler(c *fiber.Ctx) error {

	body := c.Request().Body()

	var id *ID
	err := json.Unmarshal(body, &id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "error unmarshalling request body", "ensure the request body is valid JSON containing an ID")
		return nil
	}

	pipe := h.rc.TxPipeline()

	// remove expired keys
	pipe.ZRemRangeByScore(c.Context(), id.ID, "0", strconv.FormatInt(time.Now().Add(-1*h.k.Duration("ratelimit.window")).UnixNano(), 10))

	// get existing keys
	zrange := pipe.ZRange(c.Context(), id.ID, 0, -1)

	// add new key
	pipe.ZAdd(c.Context(), id.ID, &redis.Z{Score: float64(time.Now().UnixNano()), Member: time.Now().UnixNano()})

	// set expiration
	pipe.Expire(c.Context(), id.ID, h.k.Duration("ratelimit.window"))
	_, err = pipe.Exec(c.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "error ratelimiting", "please try again later")
		return nil
	}

	// check if we're over the limit
	remaining := h.k.Int("ratelimit.limit") - len(zrange.Val())
	if remaining > 0 {
		err = c.Next()
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "error handling request", "please try again later")
			return nil
		}
	} else {
		errorResponse(c, http.StatusTooManyRequests, "ratelimit exceeded", "please try again later")
		return nil
	}
	return nil
}
