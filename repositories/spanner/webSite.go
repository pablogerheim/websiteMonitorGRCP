package spanner

import (
	"websiteMonitor/models"

	"gorm.io/gorm"
)

type SpannerResource struct {
	DB *gorm.DB
}

type ISpannerSite interface {
	CreateSite(site *models.Site) (*models.Site, error)
	GetSite(id string) (*models.Site, error)
	UpdateSite(site *models.Site) (*models.Site, error)
	DeleteSite(siteId string) error
	GetAllSites() ([]*models.Site, error)
	AutoMigrate() error
}

func (s *SpannerResource) CreateSite(site *models.Site) (*models.Site, error) {
	resp := s.DB.Create(site)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return site, nil
}

func (s *SpannerResource) GetSite(id string) (*models.Site, error) {
	var site models.Site
	resp := s.DB.First(&site, "id = ?", id)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &site, nil
}

func (s *SpannerResource) UpdateSite(site *models.Site) (*models.Site, error) {
	resp := s.DB.Save(site)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return site, nil
}

func (s *SpannerResource) DeleteSite(siteId string) error {
	var site models.Site
	resp := s.DB.Where("id = ?", siteId).Delete(&site)
	return resp.Error
}

func (s *SpannerResource) GetAllSites() ([]*models.Site, error) {
	var sites []*models.Site
	resp := s.DB.Find(&sites)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return sites, nil
}

func (s *SpannerResource) AutoMigrate() error {
	if err := s.DB.AutoMigrate(&models.Site{}); err != nil {
		return err
	}
	return nil
}

func NewSpannerRepository(DB *gorm.DB) ISpannerSite {
	return &SpannerResource{
		DB: DB,
	}
}
