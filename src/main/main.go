package main

import (
	"fmt"
	zerolog "github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	HOST = "localhost:8080"
)

func main()  {

	srv := &http.Server{
		Addr:              HOST,
		Handler:           routes(),
		ReadTimeout:       time.Second*15,
		ReadHeaderTimeout: time.Second*5,
		WriteTimeout:      time.Second*10,
		IdleTimeout:       time.Second*8,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go handlingPrettyShutdown(srv)

	<- sigCh
	return
}

func handlingPrettyShutdown(srv *http.Server)  {
	zerolog.Log().Msg(fmt.Sprintf("Server is listening on %s ...", HOST))
	if err := srv.ListenAndServe(); err != nil {
		zerolog.Fatal().Msg(err.Error())
		return
	}

	return
}