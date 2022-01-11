package facade

import (
	"gitlab.com/CalebTracey/nft-power-barn/service/nftport"
)

const (
	testChain = "rinkeby"
	testHash  = "0x47fa0ff41a8e7c1041e3392de3f7cefc676cd64ad306a78aa625537ae781eaa5"
)

type ContractFacade interface {
	GetNftPortContract() (nftport.ContractResponse, error)
	DeployNftPortContract() (nftport.ContractResponse, error)
}

type ContractService struct {
	nftport.ServiceI
}

func NewNftPortService() ContractService {
	nftPortSvc := nftport.InitializeService()
	return ContractService{
		nftPortSvc,
	}
}

func (s *ContractService) GetNftPortContract() (nftport.ContractResponse, error) {
	req := nftport.ContractRequest{
		Hash:  testHash,
		Chain: testChain,
	}
	response, err := s.ServiceI.GetContract(req)
	if err != nil {
		return nftport.ContractResponse{}, err
	}
	return response, nil
}

func (s *ContractService) DeployNftPortContract() (nftport.ContractResponse, error) {
	req := nftport.DeployContractRequest{
		Chain:        testChain,
		Name:         "SPOOKYHEADS",
		Symbol:       "SH",
		OwnerAddress: "0xE2ca23ddEF2dC6A0778bDa695517c8EB389d5e11",
	}
	response, err := s.ServiceI.DeployContract(req)
	if err != nil {
		return nftport.ContractResponse{}, err
	}

	return response, nil
}
