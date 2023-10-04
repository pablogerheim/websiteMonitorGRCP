package server

import (
	"websiteMonitor/handler"
	"websiteMonitor/pb"
)

type ServersController struct {
	Site pb.WebsiteMonitorServiceServer
}

func (s *ServersController) Inject(handler *handler.Handlers) {
	s.Site = NewSiteServer(handler)
}
