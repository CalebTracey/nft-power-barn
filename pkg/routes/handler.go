package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/CalebTracey/nft-power-barn/pkg/facade"
	"io"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	GenService     facade.GenService
	NftPortService facade.ContractFacade
	IpfsService    facade.IpfsService
}

func (h Handler) InitializeRoutes() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Handle("/generate", h.Generate()).Methods(http.MethodGet)
	r.Handle("/upload", h.UploadIpfs()).Methods(http.MethodPost)
	r.Handle("/contract", h.GetContract()).Methods(http.MethodGet)
	r.Handle("/contract", h.DeployContract()).Methods(http.MethodPost)

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

func (h Handler) UploadIpfs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := h.IpfsService.UploadCollectionIpfs()
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(result)
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

func (h Handler) DeployContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Deploying Contract...")
		//var res nftport.ContractResponse

		result, err := h.NftPortService.DeployNftPortContract()
		if err != nil {
			log.Fatalln(err)
		}

		if strings.EqualFold(result.Response, "NOK") {
			m := fmt.Sprintf("Contract Request Failed: %v", result.Error)
			buf := new(bytes.Buffer)
			err = json.NewEncoder(buf).Encode(&m)
			if err != nil {
				log.Fatalln(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(buf.Bytes())
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(&result)
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

func (h Handler) GetContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(r.Body)
		log.Println("Getting Contract...")

		result, err := h.NftPortService.GetNftPortContract()
		if err != nil {
			log.Fatalln(err)
		}
		if strings.EqualFold(result.Response, "NOK") {
			m := fmt.Sprintf("Contract Request Failed: %v", result.Error)
			buf := new(bytes.Buffer)
			err = json.NewEncoder(buf).Encode(&m)
			if err != nil {
				log.Fatalln(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(buf.Bytes())
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(&result)
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
