package routes

import (
	"myapp/controllers"
	"myapp/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	router.Handle("/items", middleware.AuthMiddleware(http.HandlerFunc(controllers.GetItems))).Methods("GET")
	router.Handle("/items", middleware.AuthMiddleware(http.HandlerFunc(controllers.CreateItem))).Methods("POST")
	router.Handle("/items/{id}", middleware.AuthMiddleware(http.HandlerFunc(controllers.DeleteItem))).Methods("DELETE")
	router.Handle("/items/{id}", middleware.AuthMiddleware(http.HandlerFunc(controllers.UpdateItem))).Methods("PUT")

}
