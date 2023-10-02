package controllers

import (
	"fmt"
	"net/http"
	"time"
	"websiteMonitor/database"
	"websiteMonitor/models"

	"github.com/gin-gonic/gin"
)

func AutoMigrate(c *gin.Context) {
	err := database.DB.AutoMigrate(&models.Site{})
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Erro ao fazer AutoMigrate",
		})
		return
	}
	c.JSON(200, gin.H{
		"success": "AutoMigrate concluído com sucesso",
	})
}

func ExibeTodosSites(c *gin.Context) {
	var sites []models.Site
	var response []models.SiteResponse

	database.DB.Find(&sites)

	for _, site := range sites {
		resp, err := http.Get(site.Name)
		if err != nil {
			fmt.Println("Ocorreu um erro:", err)
			continue
		}

		siteResponse := models.SiteResponse{
			ID:   int(site.ID),
			Name: site.Name,
			Date: time.Now().Format(time.RFC3339),
		}

		if resp.StatusCode < 400 {
			siteResponse.Status = "Online"
		} else {
			siteResponse.Status = "Offline"
		}

		response = append(response, siteResponse)
	}
	fmt.Println(response)

	c.JSON(http.StatusOK, response)
}

func CriaNovoSite(c *gin.Context) {
	var site models.Site
	if err := c.ShouldBindJSON(&site); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	database.DB.Create(&site)

	c.JSON(http.StatusOK, site)
}

func DeletaSite(c *gin.Context) {
	var site models.Site
	id := c.Params.ByName("id")
	database.DB.Delete(&site, id)
	c.JSON(http.StatusOK, gin.H{"data": "Site deletado com sucesso"})
}

func EditaSite(c *gin.Context) {
	var site models.Site
	id := c.Params.ByName("id")
	database.DB.First(&site, id)

	if err := c.ShouldBindJSON(&site); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	database.DB.Model(&site).UpdateColumns(site)
	c.JSON(http.StatusOK, site)
}
