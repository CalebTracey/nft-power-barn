package generate

import (
	"gitlab.com/CalebTracey/nft-power-barn/pkg/facade"
)

type ServiceI interface {
	Start()
}

type Service struct {
	Service facade.GenFacade
}

func (s Service) Start() {

	s.Service.StartCreating()
}
