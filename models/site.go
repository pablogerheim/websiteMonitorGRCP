package models

import "gorm.io/gorm"

type Site struct {
	gorm.Model
	Name string `json:"name"`
}

type SiteResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Date   string `json:"date"`
	Status string `json:"status"`
}
