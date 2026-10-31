package rates

import (
	"net/http"

	"github.com/google/uuid"

	"exchanger/internal/ports/http/response"
)

// getUpdate godoc
// @Summary  Get an update by id
// @Tags     rates
// @Produce  json
// @Param    id   path      string  true  "Update id"  format(uuid)
// @Success  200  {object}  updateResponse
// @Failure  400  {object}  response.ErrorResponse
// @Failure  404  {object}  response.ErrorResponse
// @Failure  500  {object}  response.ErrorResponse
// @Router   /api/v1/rates/updates/{id} [get]
func (c *Controller) getUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid update id")
		return
	}

	update, err := c.service.GetUpdate(r.Context(), id)
	if err != nil {
		c.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, newUpdateResponse(update))
}
