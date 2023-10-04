package handler

import "websiteMonitor/repositories"

type Handlers struct {
	Site ISiteHandler
}

func (h *Handlers) Inject(repo *repositories.Repositories) {
	h.Site = NewSiteHandler(repo)
}
