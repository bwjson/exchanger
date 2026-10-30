package rates

import (
	"net/http"

	"github.com/google/uuid"

	"exchanger/internal/ports/http/response"
)

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
