package rates

import (
	"encoding/json"
	"net/http"

	"exchanger/internal/ports/http/response"
)

const maxIdempotencyKeyLength = 255

// createUpdate godoc
// @Summary  Request a rate update
// @Description  Creates an update for a pair and returns its id. The rate is fetched in the background.
// @Tags     rates
// @Accept   json
// @Produce  json
// @Param    Idempotency-Key  header    string               false  "A repeated key returns the same update_id"  maxlength(255)
// @Param    request          body      createUpdateRequest  true   "Currency pair"
// @Success  202              {object}  createUpdateResponse
// @Failure  400              {object}  response.ErrorResponse
// @Failure  422              {object}  response.ErrorResponse
// @Failure  500              {object}  response.ErrorResponse
// @Router   /api/v1/rates/updates [post]
func (c *Controller) createUpdate(w http.ResponseWriter, r *http.Request) {
	var req createUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	key := r.Header.Get("Idempotency-Key")
	if len(key) > maxIdempotencyKeyLength {
		response.Error(w, http.StatusBadRequest, "idempotency key is too long")
		return
	}

	id, err := c.service.CreateUpdate(r.Context(), req.Pair, key)
	if err != nil {
		c.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusAccepted, createUpdateResponse{UpdateID: id})
}
