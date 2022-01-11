package nftport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	DeploycontracturlTest = "https://api.nftport.xyz/v0/contracts/transaction_hash?chain=rinkeby"
	GetcontracturlTest    = "https://api.nftport.xyz/v0/contracts"
)

type ServiceI interface {
	GetContract(req ContractRequest) (ContractResponse, error)
	DeployContract(req DeployContractRequest) (ContractResponse, error)
}

type Service struct {
	GetContractUrl    string
	DeployContractUrl string
	Client            *http.Client
}

func InitializeService() *Service {
	getConUrl := GetcontracturlTest
	DepConUrl := DeploycontracturlTest
	c := &http.Client{}

	return &Service{
		GetContractUrl:    getConUrl,
		DeployContractUrl: DepConUrl,
		Client:            c,
	}
}

func (s *Service) GetContract(contractReq ContractRequest) (ContractResponse, error) {
	var res ContractResponse
	//urlParams := fmt.Sprintf("%v%v?chain=%v", s.GetContractUrl, contractReq.Hash, contractReq.Chain)
	req, err := http.NewRequest("GET", "https://api.nftport.xyz/v0/contracts/0x47fa0ff41a8e7c1041e3392de3f7cefc676cd64ad306a78aa625537ae781eaa5?chain=rinkeby", nil)
	if err != nil {
		return ContractResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "6d68a19d-3f75-442f-a593-866763c9d6e2")

	r, err := s.Client.Do(req)
	if err != nil {
		return ContractResponse{}, err
	}
	if r.StatusCode != http.StatusOK {
		return ContractResponse{}, fmt.Errorf("failed to fetch contract")
	}
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return ContractResponse{}, err
	}

	log.Println(res)
	return res, nil
}

func (s *Service) DeployContract(payload DeployContractRequest) (ContractResponse, error) {
	var res ContractResponse
	buff := new(bytes.Buffer)
	err := json.NewEncoder(buff).Encode(payload)

	if err != nil {
		return ContractResponse{}, err
	}
	req, err := http.NewRequest("POST", s.DeployContractUrl, buff)
	if err != nil {
		return ContractResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "6d68a19d-3f75-442f-a593-866763c9d6e2")

	r, err := s.Client.Do(req)
	if err != nil {
		return ContractResponse{}, err
	}

	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return ContractResponse{}, err
	}

	return res, nil
}
