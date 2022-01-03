package main

import (
	"generatecollection/pkg/facade"
	"generatecollection/pkg/services/generate"
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
