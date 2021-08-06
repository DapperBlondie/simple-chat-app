package main

import (
	"github.com/DapperBlondie/simple-chat-app/src/handlers"
	"github.com/go-chi/chi"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/home", handlers.AppConf.Home)
	mux.Get("/chat-app", handlers.AppConf.WsEndpointHandler)

	return mux
}
