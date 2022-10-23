package handlers

import (
	"encoding/json"
	"net/http"

	"676f.dev/rfc7807"
	"github.com/gofiber/fiber/v2"
)

type ID struct {
	ID string `json:"id"`
}

func (h *Handler) IDHandler(c *fiber.Ctx) error {
	body := c.Request().Body()

	var id *ID
	err := json.Unmarshal(body, &id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "error unmarshalling request body", "ensure the request body is valid JSON containing an ID")
		return nil
	}

	respBody := rfc7807.NewError("id handler").
		SetStatus(http.StatusOK).
		SetExtensions(map[string]any{
			"id": id.ID,
		})

	bytes, err := json.Marshal(respBody)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "error marshalling response body", "please try again later")
		return nil
	}

	c.Status(http.StatusOK)
	_, err = c.Write(bytes)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "error writing response body", "please try again later")
		return nil
	}

	return nil
}
