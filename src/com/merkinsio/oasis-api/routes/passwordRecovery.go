package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendPasswordRecoveryRoutes(router *mux.Router) {
	router.Handle("/", server.UserRecoveryMiddleware(server.PasswordRecovery)).Methods("POST")
	router.Handle("/token/{token:[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-4[0-9A-Fa-f]{3}-[89ABab][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12}}", server.UserRecoveryMiddleware(server.ChangePasswordRecovery)).Methods("POST")
}
