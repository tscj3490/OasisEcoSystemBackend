package routes

import (
	"com/merkinsio/oasis-api/server"

	"github.com/gorilla/mux"
)

// AppendIssuesAPI Holds all the /issues routes
func AppendIssuesAPI(router *mux.Router) {
	router.Handle("/", server.AuthenticateWithUser(server.CreateIssue)).Methods("POST")
	router.Handle("/", server.AuthenticateWithUser(server.GetAllIssues)).Methods("GET")
	//TODO: ADD "/?year=yyyy" method handler
	router.Handle("/{issueId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.UpdateIssue)).Methods("POST")
	router.Handle("/{issueId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.GetIssueByID)).Methods("GET")
	router.Handle("/{issueId:[a-f\\d]{24}}/work", server.AuthenticateWithUser(server.StartWorkingIssue)).Methods("POST")
	router.Handle("/{issueId:[a-f\\d]{24}}/pause", server.AuthenticateWithUser(server.PauseWorkingIssue)).Methods("POST")
	router.Handle("/{issueId:[a-f\\d]{24}}/finish", server.AuthenticateWithUser(server.FinishIssue)).Methods("POST")
	router.Handle("/{issueId:[a-f\\d]{24}}/accept", server.AuthenticateWithUser(server.AcceptWorkOrderIssue)).Methods("POST")
	router.Handle("/finished", server.AuthenticateWithUser(server.GetAllFinishedIssues)).Methods("GET")
	router.Handle("/{issueId:[a-f\\d]{24}}", server.AuthenticateWithUser(server.DeleteIssue)).Methods("DELETE")
	router.Handle("/task", server.AuthenticateWithUser(server.CreateTaskIssue)).Methods("POST")
	router.Handle("/states", server.AuthenticateWithUser(server.GetIssueStatuses)).Methods("GET")
	router.Handle("/treatments", server.AuthenticateWithUser(server.GetTreatmentTypes)).Methods("GET")
}
