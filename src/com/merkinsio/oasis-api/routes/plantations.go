package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

// AppendPlantationsAPI Holds all the /plantations routes
func AppendPlantationsAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.CreatePlantation)).Methods("POST")
	router.Handle("/{plantationId:[A-Fa-f\\d]{24}}", server.AuthenticateWithUser(server.UpdatePlantationByID)).Methods("POST")
	router.Handle("/{plantationId:[A-Fa-f\\d]{24}}", server.AuthenticateWithUser(server.DeletePlantation)).Methods("DELETE")
	router.Handle("/", server.AuthenticateWithUser(server.GetPlantationsByCurrentUser)).Methods("GET")
	router.Handle("/plots", server.AuthenticateWithUser(server.GetPlantationsAndPlotsByCurrentUser)).Methods("GET")
	router.Handle("/{plantationId:[A-Fa-f\\d]{24}}/issues", server.AuthenticateWithUser(server.GetPlantationIssues)).Methods("GET")
	router.Handle("/{plantationId:[A-Fa-f\\d]{24}}/issues/finished", server.AuthenticateWithUser(server.GetPlantationFinishedIssues)).Methods("GET")
	router.Handle("/{plantationId:[A-Fa-f\\d]{24}}/sigpacPlots", server.AuthenticateWithUser(server.GetPlantationSigpacPlots)).Methods("GET")
	router.HandleFunc("/{plantationId:[A-Fa-f\\d]{24}}/sigpacPlotsWithNoUser", server.GetPlantationSigpacPlotsWithNoUser).Methods("GET")
	router.Handle("/{plantationId:[A-Fa-f\\d]{24}}/myIssues", server.AuthenticateWithUser(server.GetMyIssues)).Methods("GET")
}
