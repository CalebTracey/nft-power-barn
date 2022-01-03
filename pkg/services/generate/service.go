package generate

import (
	"generatecollection/pkg/facade"
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
