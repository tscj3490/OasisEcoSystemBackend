package main

import (
	"com/merkinsio/oasis-api/config"
	"com/merkinsio/oasis-api/domain"
	"com/merkinsio/oasis-api/routes"
	"com/merkinsio/oasis-api/server"
	"net"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

var jwtMiddleware *jwtmiddleware.JWTMiddleware

// App Structure that holds the application
type App struct {
	Router *mux.Router
	DB     *domain.MongoConnector
	// Context  *server.Context
	Listener net.Listener
	// SubRouters map[string]*mux.Route
}

// Initialize This function initialize the application router and context
func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	// a.Context = server.InitContext()

	a.DB = domain.NewConnector(config.Config.GetString("database.url"))

	jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("9cdd3cbd-b48f-4187-a48a-07e67e2e9ab3"), nil //TODO: CHANGE MY SECRET
		},
		Extractor: jwtmiddleware.FromFirst(
			jwtmiddleware.FromAuthHeader,
			jwtmiddleware.FromParameter("auth_token"),
		),
		SigningMethod: jwt.SigningMethodHS256,
	})

	a.initializeRoutes()
}

// ShutDown This function cleans the context in the application
func (a *App) ShutDown() {
	if a.Listener != nil {
		a.Listener.Close()
	}
	a.DB.Close()
}

// Run This function starts the application at the designed port
func (a *App) Run(addr string) {
	listener, err := net.Listen("tcp", addr)
	debug := false
	loggerLevel := logrus.ErrorLevel

	if config.Config.GetString("environment") == "development" {
		logrus.Info("Using development mode")
		debug = true
		loggerLevel = logrus.DebugLevel
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		Debug:            debug,
	})

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetLevel(loggerLevel)

	n := negroni.New()
	// headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	// originsOk := handlers.AllowedOrigins([]string{"*"})
	// methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Error handler
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false

	n.Use(c)
	n.Use(negronilogrus.NewMiddlewareFromLogger(logger, "oasis-api"))
	n.Use(recovery)
	n.UseHandler(a.Router)

	if err != nil {
		panic(err.Error())
	}
	a.Listener = listener
	// log.Fatal(http.ListenAndServe(addr, a.Router))
	logger.Fatal(http.Serve(a.Listener, n))
}

func (a *App) initializeRoutes() {

	// Define the main api path
	apiRouter := a.Router.PathPrefix("/api").Subrouter()

	// Health check
	apiRouter.HandleFunc("/health-check", server.HealthCheckHandler).Methods("GET")

	// Websocket endpoint
	apiRouter.Handle("/ws", server.AuthenticateWithUser(server.WebsocketHandler))

	// PasswordRecovery endpoint
	passwordRecoveryRouter := apiRouter.PathPrefix("/passwordRecovery").Subrouter().StrictSlash(true)
	routes.AppendPasswordRecoveryRoutes(passwordRecoveryRouter)

	authenticateRouter := apiRouter.PathPrefix("/authenticate").Subrouter().StrictSlash(true)
	routes.AppendAuthenticateRoutes(authenticateRouter)
	apiRouter.Handle("/me", server.AuthenticateWithUser(server.Me)).Methods("GET")
	apiRouter.Handle("/me/{versionCode:[0-9]+}", server.AuthenticateWithUser(server.Me)).Methods("GET")

	/* Domain routes */
	// Here goes the other routes
	usersRouter := apiRouter.PathPrefix("/users").Subrouter().StrictSlash(true)
	routes.AppendUsersAPI(usersRouter)

	plotsRouter := apiRouter.PathPrefix("/plots").Subrouter().StrictSlash(true)
	routes.AppendPlotsAPI(plotsRouter)

	plantationsRouter := apiRouter.PathPrefix("/plantations").Subrouter().StrictSlash(true)
	routes.AppendPlantationsAPI(plantationsRouter)

	issuesRouter := apiRouter.PathPrefix("/issues").Subrouter().StrictSlash(true)
	routes.AppendIssuesAPI(issuesRouter)

	filesRouter := apiRouter.PathPrefix("/files").Subrouter().StrictSlash(true)
	routes.AppendFilesRouter(filesRouter)

	cooperativesRouter := apiRouter.PathPrefix("/cooperatives").Subrouter().StrictSlash(true)
	routes.AppendCooperativesAPI(cooperativesRouter)

	visitCardRouter := apiRouter.PathPrefix("/visitCards").Subrouter().StrictSlash(true)
	routes.AppendVisitCardRoutes(visitCardRouter)

	paramsRouter := apiRouter.PathPrefix("/params").Subrouter().StrictSlash(true)
	routes.AppendParamsAPI(paramsRouter)

	productsRouter := apiRouter.PathPrefix("/products").Subrouter().StrictSlash(true)
	routes.AppendProductsAPI(productsRouter)

	employeesRouter := apiRouter.PathPrefix("/employees").Subrouter().StrictSlash(true)
	routes.AppendEmployeesAPI(employeesRouter)

	statisticsRouter := apiRouter.PathPrefix("/statistics").Subrouter().StrictSlash(true)
	routes.AppendStatisticsAPI(statisticsRouter)

	sigpacCropsRouter := apiRouter.PathPrefix("/sigpacCrops").Subrouter().StrictSlash(true)
	routes.AppendSigpacCropsAPI(sigpacCropsRouter)

	declaredCultivationCropsRouter := apiRouter.PathPrefix("/declaredCultivationsCrops").Subrouter().StrictSlash(true)
	routes.AppendDeclaredCultivationCropsAPI(declaredCultivationCropsRouter)

	variatyCropsRouter := apiRouter.PathPrefix("/variatyCrops").Subrouter().StrictSlash(true)
	routes.AppendVariatyCropsAPI(variatyCropsRouter)

	sigpacPlotsRouter := apiRouter.PathPrefix("/sigpacPlots").Subrouter().StrictSlash(true)
	routes.AppendSigpacPlotsAPI(sigpacPlotsRouter)

	ProvincesRouter := apiRouter.PathPrefix("/provinces").Subrouter().StrictSlash(true)
	routes.AppendProvincesAPI(ProvincesRouter)

	TownsRouter := apiRouter.PathPrefix("/towns").Subrouter().StrictSlash(true)
	routes.AppendTownsAPI(TownsRouter)
	/* End*/

}
