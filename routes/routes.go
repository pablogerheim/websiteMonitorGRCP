package routes

import (
	"websiteMonitor/controllers"

	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()
	r.GET("/autoMigrate", controllers.AutoMigrate)
	r.GET("/site", controllers.ExibeTodosSites)
	r.POST("/site", controllers.CriaNovoSite)
	r.DELETE("/site/:id", controllers.DeletaSite)
	r.PATCH("/site/:id", controllers.EditaSite)
	r.GET("/iniciar", controllers.IniciarRotina)
	r.GET("/parar", controllers.PararRotina)
	r.Run()
}
