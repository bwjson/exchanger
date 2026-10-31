package rates

import (
	"net/http"

	"exchanger/internal/ports/http/response"
)

// getLatest godoc
// @Summary  Get the last known rate of a pair
// @Tags     rates
// @Produce  json
// @Param    pair  query     string  true  "Currency pair"  example(EUR/MXN)
// @Success  200   {object}  rateResponse
// @Failure  400   {object}  response.ErrorResponse
// @Failure  404   {object}  response.ErrorResponse
// @Failure  500   {object}  response.ErrorResponse
// @Router   /api/v1/rates/latest [get]
func (c *Controller) getLatest(w http.ResponseWriter, r *http.Request) {
	rate, err := c.service.GetLatest(r.Context(), r.URL.Query().Get("pair"))
	if err != nil {
		c.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, newRateResponse(rate))
}
