package worker

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Api struct {
	Address string
	Port    int
	Worker  *Worker
	Router  *chi.Mux
}

type ErrResponse struct {
	HTTPStatusCode int
	Message        string
}

func (api *Api) initRouter() {
	api.Router = chi.NewRouter()
	api.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", api.StartTaskHandler)
		r.Get("/", api.GetTaskHandler)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", api.StopTaskHandler)
		})
	})
}

func (api *Api) Start() {
	api.initRouter()
	url := fmt.Sprintf("%s:%d", api.Address, api.Port)
	http.ListenAndServe(url, api.Router)
}
