package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"676f.dev/rfc7807"
	"github.com/go-redis/redis/v8"
	"github.com/knadh/koanf"
)

type IDHandler struct {
	rc *redis.Client
	k  *koanf.Koanf
}

type ID struct {
	ID string `json:"id"`
}

func NewIDHandler(rc *redis.Client, k *koanf.Koanf) *IDHandler {
	return &IDHandler{
		rc: rc,
		k:  k,
	}
}

func (idh *IDHandler) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = rfc7807.SimpleResponse(w, http.StatusInternalServerError, "error reading request body")
		return
	}

	var id *ID
	err = json.Unmarshal(body, &id)
	if err != nil {
		_ = rfc7807.SimpleResponse(w, http.StatusInternalServerError, "error unmarshalling request body")
		return
	}

	isAllowed, err := idh.ratelimit(r.Context(), id.ID)
	if err != nil {
		_ = rfc7807.SimpleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !isAllowed {
		_ = rfc7807.SimpleResponse(w, http.StatusTooManyRequests, "rate limit exceeded")
		return
	}

	respBody := rfc7807.NewError("id handler").
		SetStatus(http.StatusOK).
		SetExtensions(map[string]any{
			"id": id.ID,
		})

	bytes, err := json.Marshal(respBody)
	if err != nil {
		_ = rfc7807.SimpleResponse(w, http.StatusInternalServerError, "error marshalling response body")
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		_ = rfc7807.SimpleResponse(w, http.StatusInternalServerError, "error writing response body")
		return
	}
}

func (idh *IDHandler) ratelimit(ctx context.Context, id string) (bool, error) {
	pipe := idh.rc.TxPipeline()

	// remove expired keys
	pipe.ZRemRangeByScore(ctx, id, "0", strconv.FormatInt(time.Now().Add(-1*idh.k.Duration("ratelimit.window")).UnixNano(), 10))

	// get existing keys
	zrange := pipe.ZRange(ctx, id, 0, -1)

	// add new key
	pipe.ZAdd(ctx, id, &redis.Z{Score: float64(time.Now().UnixNano()), Member: time.Now().UnixNano()})

	// set expiration
	pipe.Expire(ctx, id, idh.k.Duration("ratelimit.window"))
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	// check if we're over the limit
	remaining := idh.k.Int("ratelimit.limit") - len(zrange.Val())
	if remaining > 0 {
		return true, nil
	}
	return false, nil
}
