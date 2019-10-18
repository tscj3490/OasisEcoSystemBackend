package routes

import (
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

// AppendUsersAPI Holds all the /users routes
func AppendUsersAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.GetAllUsers)).Methods("GET")
	router.Handle("/", server.AuthenticateWithUser(server.CreateUser)).Methods("POST")
	router.Handle("/{userId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.GetUserByID)).Methods("GET")
	router.Handle("/{userId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.UpdateUser)).Methods("POST")
	router.Handle("/{userId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.DeleteUser)).Methods("DELETE")
	router.Handle("/{userId:[a-f\\d]{24}}/termofservices", server.AuthenticateWithUser(server.TermOfServicesAcceptedByUser)).Methods("POST")
	router.Handle("/register", server.Wrap(server.RegisterUser)).Methods("POST")
	router.Handle("/technicians", server.AuthenticateWithUser(server.GetAllTechnicians)).Methods("GET")
	router.Handle("/workers", server.AuthenticateWithUser(server.GetAllWorkers)).Methods("GET")
	router.Handle("/cooperative/{cooperativeId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.UserMatchRoles(constants.RoleTechnician, constants.RoleAdmin, constants.RoleWorker), server.GetClientsByCooperativeID)).Methods("GET")
}
