package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

// AppendPlotsAPI Holds all the /plots routes
func AppendPlotsAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.CreatePlot)).Methods("POST")
	router.Handle("/{plotId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.UpdatePlot)).Methods("POST")
	router.Handle("/{plotId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.DeletePlot)).Methods("DELETE")
	router.Handle("/{plotId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.GetPlotByID)).Methods("GET")
	router.Handle("/{plotId:[a-f\\d]{24}}/sigpac", server.AuthenticateWithUser(server.GetSigPacPlotByPlotID)).Methods("GET")
	router.Handle("/plantation/{plantationId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.GetPlotsByPlantationID)).Methods("GET")
	router.Handle("/plantation/{plantationId:[a-f\\d]{24}}/geo", server.AuthenticateWithUser(server.GetPlotsWithGeoJSONByPlantationID)).Methods("GET")
	router.HandleFunc("/{plotId:[a-f\\d]{24}}/sigpacPlotWithNoUser", server.GetSigPacPlotByPlotIDWithNoUser).Methods("GET")
}
