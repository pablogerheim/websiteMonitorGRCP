package handler

import (
	"websiteMonitor/models"
	"websiteMonitor/repositories"
)

type resource struct {
	repository *repositories.Repositories
}

type ISiteHandler interface {
	CreateSite(site *models.Site) (*models.Site, error)
	GetSite(id string) (*models.Site, error)
	UpdateSite(site *models.Site) (*models.Site, error)
	DeleteSite(siteID string) error
	GetAllSites() ([]*models.Site, error)
	AutoMigrate() error
}

func (r *resource) CreateSite(site *models.Site) (*models.Site, error) {
	resp, err := r.repository.Spanner.CreateSite(site)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resource) GetSite(id string) (*models.Site, error) {
	resp, err := r.repository.Spanner.GetSite(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resource) UpdateSite(site *models.Site) (*models.Site, error) {
	resp, err := r.repository.Spanner.UpdateSite(site)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *resource) DeleteSite(siteID string) error {
	err := r.repository.Spanner.DeleteSite(siteID)
	if err != nil {
		return err
	}
	return nil
}

func (r *resource) GetAllSites() ([]*models.Site, error) {
	sites, err := r.repository.Spanner.GetAllSites()
	if err != nil {
		return nil, err
	}
	return sites, nil
}

func (r *resource) AutoMigrate() error {
	err := r.repository.Spanner.AutoMigrate()
	if err != nil {
		return err
	}
	return nil
}

func NewSiteHandler(hspanner *repositories.Repositories) ISiteHandler {
	return &resource{
		repository: hspanner,
	}
}
