package workbench

import (
	"embed"
	"net/http"
)

//go:embed static/*
var staticFS embed.FS

type Server struct {
	Handler http.Handler
}

func New(api http.Handler) *Server {
	mux := http.NewServeMux()
	mux.Handle("/api/", api)
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	return &Server{Handler: mux}
}
