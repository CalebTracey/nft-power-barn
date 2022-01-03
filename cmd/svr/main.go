package main

import (
	//"generatecollection/pkg/facade"
	//"generatecollection/pkg/facade"
	"gitlab.com/CalebTracey/nft-power-barn/pkg/facade"
	"gitlab.com/CalebTracey/nft-power-barn/pkg/services/generate"
	"log"
)

func main() {
	defer deathScream()

	genService := facade.NewService()
	gen := generate.Service{
		Service: &genService,
	}
	gen.Start()
}

func deathScream() {
	if r := recover(); r != nil {
		log.Panicf("I panicked and am quitting: %v", r)
	}
}
