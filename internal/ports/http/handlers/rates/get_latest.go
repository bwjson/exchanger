package rates

import (
	"net/http"

	"exchanger/internal/ports/http/response"
)

func (c *Controller) getLatest(w http.ResponseWriter, r *http.Request) {
	pair := r.URL.Query().Get("pair")
	if pair == "" {
		response.Error(w, http.StatusBadRequest, "pair is required")
		return
	}

	response.Error(w, http.StatusNotImplemented, "not implemented")
}
