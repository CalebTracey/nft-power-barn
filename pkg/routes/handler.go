package routes

import (
	"github.com/gorilla/mux"
	"gitlab.com/CalebTracey/nft-power-barn/pkg/facade"
	"net/http"
)

type Handler struct {
	Service facade.GenService
}

func (h Handler) InitializeRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/generate", h.Generate()).Methods(http.MethodGet)
	r.Handle("/nftport", h.DeployContract()).Methods(http.MethodPost)

	return r
}

func (h Handler) Generate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Service.StartCreating()
	}
}

func (h Handler) DeployContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
