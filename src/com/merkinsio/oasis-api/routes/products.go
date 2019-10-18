package routes

import (
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendProductsAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.GetAllProducts)).Methods("GET")
	router.Handle("/", server.AuthenticateWithUser(server.UserMatchRoles(constants.RoleAdmin), server.CreateProduct)).Methods("POST")
	router.Handle("/{productId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.UserMatchRoles(constants.RoleAdmin), server.DeleteProduct)).Methods("DELETE")
}
