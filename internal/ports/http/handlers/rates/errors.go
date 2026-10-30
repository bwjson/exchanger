package rates

import (
	"errors"
	"net/http"

	"exchanger/internal/domain"
	"exchanger/internal/ports/http/response"
)

func (c *Controller) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidPair):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	default:
		c.logger.Error("request failed", "error", err)
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
