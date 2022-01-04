package main

import (
	"context"
	"flag"
	"gitlab.com/CalebTracey/nft-power-barn/pkg/facade"
	"os"
	"os/signal"

	//"generatecollection/pkg/facade"
	//"generatecollection/pkg/facade"
	"gitlab.com/CalebTracey/nft-power-barn/pkg/routes"
	"log"
	"net/http"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	defer deathScream()

	service := facade.NewService()
	handler := routes.Handler{
		Service: service,
	}

	router := handler.InitializeRoutes()

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	_ = srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

func deathScream() {
	if r := recover(); r != nil {
		log.Panicf("I panicked and am quitting: %v", r)
	}
}
