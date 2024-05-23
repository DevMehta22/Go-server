package routes

import (
	"github.com/DevMehta22/mongoapi/controller"
	"github.com/gorilla/mux"
)

func Router()  *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/movies",controller.GetAllMovies).Methods("GET")
	router.HandleFunc("/api/movie",controller.CreateMovie).Methods("POST")
	router.HandleFunc("/api/movie/{id}",controller.MarkAsWatched).Methods("PUT")
	router.HandleFunc("/api/movie/{id}",controller.DeleteMove).Methods("DELETE")
	router.HandleFunc("/api/deleteallmovie",controller.DeleteMovies).Methods("DELETE")

	return router
}