package routes

import "gitlab.com/CalebTracey/nft-power-barn/pkg/facade"

type Handler struct {
	Service facade.GenService
}
