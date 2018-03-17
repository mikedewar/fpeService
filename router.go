package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *FPEServer) NewRouter() *mux.Router {
	type Route struct {
		Name        string
		Pattern     string
		Method      string
		HandlerFunc http.HandlerFunc
	}
	routes := []Route{
		Route{
			"Encyrpt",
			"/encrypt",
			"GET",
			s.encryptHandler,
		},
	}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}
