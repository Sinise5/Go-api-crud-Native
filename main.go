package main

import (
	"log"
	"myapp/config"
	"myapp/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.InitDB()

	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
