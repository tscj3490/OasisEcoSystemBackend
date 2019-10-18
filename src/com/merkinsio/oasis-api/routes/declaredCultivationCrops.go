package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendDeclaredCultivationCropsAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.GetAllDeclaredCultivationCrops)).Methods("GET")
	router.Handle("/", server.AuthenticateWithUser(server.CreateDeclaredCultivationCrop)).Methods("POST")
	router.Handle("/{declaredCultivationId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.DeleteDeclaredCultivationCrop)).Methods("DELETE")
	router.Handle("/{declaredCultivationId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.UpdateDeclaredCultivationCrop)).Methods("POST")
}
