package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendFilesRouter(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.UploadFile)).Methods("POST")
	router.Handle("/{fileId:[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-4[0-9A-Fa-f]{3}-[89ABab][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12}}", server.Wrap(server.DownloadFile)).Methods("GET")
}
