package rates

import (
	"net/http"

	"exchanger/internal/ports/http/response"
)

func (c *Controller) getLatest(w http.ResponseWriter, r *http.Request) {
	rate, err := c.service.GetLatest(r.Context(), r.URL.Query().Get("pair"))
	if err != nil {
		c.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, newRateResponse(rate))
}
