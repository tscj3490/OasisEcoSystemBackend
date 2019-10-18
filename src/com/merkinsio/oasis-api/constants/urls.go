package constants

import (
	"com/merkinsio/oasis-api/config"
)

//BaseURL Holds the API url
var BaseURL = config.Config.GetString("server.baseUrl")
