package routes

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/CalebTracey/nft-power-barn/pkg/facade"
	"log"
	"net/http"
)

type Handler struct {
	GenService facade.GenService
}

func (h Handler) InitializeRoutes() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Handle("/generate", h.Generate()).Methods(http.MethodGet)

	return r
}

func (h Handler) Generate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := h.GenService.StartCreating()
		if err != nil {
			log.Fatalln(err)
		}
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(result)
		if err != nil {
			log.Fatalln(err)
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(buf.Bytes())
		if err != nil {
			log.Fatalln(err)
		}
	}
}
