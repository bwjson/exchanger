package rates

import (
	"net/http"

	"exchanger/internal/ports/http/response"
)

func (c *Controller) getUpdate(w http.ResponseWriter, r *http.Request) {
	_ = r.PathValue("id")

	response.Error(w, http.StatusNotImplemented, "not implemented")
}
