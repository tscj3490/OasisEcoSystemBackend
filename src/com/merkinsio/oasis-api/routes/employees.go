package routes

import (
	"com/merkinsio/oasis-api/constants"
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

func AppendEmployeesAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.UserMatchRoles(constants.RoleTechnician, constants.RoleAdmin, constants.RoleWorker), server.GetAllEmployees)).Methods("GET")
	router.Handle("/cooperative/{cooperativeId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.UserMatchRoles(constants.RoleTechnician, constants.RoleAdmin, constants.RoleWorker), server.GetEmployeesByCooperativeID))
}
