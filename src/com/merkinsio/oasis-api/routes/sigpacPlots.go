package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendSigpacPlotsAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.GetSigpacPlot)).Queries("provinceRec", "{^[0-9]+$}").Queries("townRec", "{townRec}").Queries("polygonNum", "{polygonNum}").Queries("plotNum", "{plotNum}").Queries("enclosure", "{enclosure}").Methods("GET")
}
