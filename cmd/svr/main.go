package main

import (
	"context"
	"flag"
	"fmt"
	"gitlab.com/CalebTracey/nft-power-barn/pkg/facade"
	"os"
	"os/signal"

	"gitlab.com/CalebTracey/nft-power-barn/pkg/routes"
	"log"
	"net/http"
	"time"
)

const Port = 8080

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	defer deathScream()

	gSvc := facade.NewGenService()
	nftSvc := facade.NewNftPortService()
	ipfsSvc := facade.NewIpfsService()
	handler := routes.Handler{
		GenService:     gSvc,
		NftPortService: &nftSvc,
		IpfsService:    ipfsSvc,
	}
	router := handler.InitializeRoutes()

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", Port),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 120,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	msg := fmt.Sprintf("Listening on Port: %v", Port)
	fmt.Println("\033[36m", msg, "\033[0m")

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
