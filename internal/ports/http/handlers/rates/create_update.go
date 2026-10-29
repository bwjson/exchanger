package rates

import (
	"encoding/json"
	"net/http"

	"exchanger/internal/ports/http/response"
)

func (c *Controller) createUpdate(w http.ResponseWriter, r *http.Request) {
	var req createUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json payload")
		return
	}
	if req.Pair == "" {
		response.Error(w, http.StatusBadRequest, "pair is required")
		return
	}

	response.Error(w, http.StatusNotImplemented, "not implemented")
}
