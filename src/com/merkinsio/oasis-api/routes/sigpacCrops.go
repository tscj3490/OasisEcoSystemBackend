package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendSigpacCropsAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.GetAllSigpacCrops)).Methods("GET")
	router.Handle("/{sigpacCropId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.GetSigpacCropById)).Methods("GET")
}
