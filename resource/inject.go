package resource

import (
	"websiteMonitor/handler"
	"websiteMonitor/repositories"
	"websiteMonitor/server"
)

var Repositories repositories.Repositories
var Handlers handler.Handlers
var ServersControllers server.ServersController

func Inject() {

	spannerClient := ConectaComBancoDeDados()

	Repositories.Inject(spannerClient)

	Handlers.Inject(&Repositories)

	ServersControllers.Inject(&Handlers)
}
