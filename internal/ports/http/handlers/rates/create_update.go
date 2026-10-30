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

	id, err := c.service.CreateUpdate(r.Context(), req.Pair)
	if err != nil {
		c.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusAccepted, createUpdateResponse{UpdateID: id})
}
