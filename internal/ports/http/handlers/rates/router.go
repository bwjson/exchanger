package rates

import "net/http"

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Init(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/rates/updates", c.createUpdate)
	mux.HandleFunc("GET /api/v1/rates/updates/{id}", c.getUpdate)
	mux.HandleFunc("GET /api/v1/rates/latest", c.getLatest)
}
