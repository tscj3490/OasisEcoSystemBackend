package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendStatisticsAPI(router *mux.Router) {
	router.Handle("/{plantationId:[A-Fa-f\\d]{24}}", server.AuthenticateWithUser(server.GetPlantationStatistics)).Methods("GET")
}
