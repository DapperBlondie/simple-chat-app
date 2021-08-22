package main

import (
	"github.com/DapperBlondie/simple-chat-app/src/handlers"
	"github.com/go-chi/chi"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("./static"))

	mux.Get("/home", handlers.AppConf.Home)
	mux.Get("/chat-app", handlers.AppConf.WsEndpointHandler)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
