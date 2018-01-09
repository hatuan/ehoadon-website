// routers.go

package routers

import (
	"erpvietnam/ehoadon-website/controllers"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()

	router.LoadHTMLGlob("templates/*")

	// Initialize the routes
	initializeRoutes()

	// Static files
	router.Static("/assets", "./assets")
	router.Static("/content", "./content")
	router.Static("/plugins", "./plugins")
	router.Static("/images", "./images")
	router.Static("/scripts", "./scripts")
	router.Static("/js", "./scripts")
	router.Static("/styles", "./styles")
	router.Static("/css", "./styles")
	router.StaticFile("/favicon.ico", "./favicon.ico")
}

func GetRoute() *gin.Engine {
	return router
}

func initializeRoutes() {

	router.GET("/", controllers.ShowIndexPage)
	router.GET("/index.html", controllers.ShowIndexPage)
	router.GET("/index", controllers.ShowIndexPage)

	router.GET("/register", controllers.ShowRegisterPage)
	router.GET("/register.html", controllers.ShowRegisterPage)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	invoiceRoutes := router.Group("/invoices")
	{
		// Handle GET requests at /invoices
		invoiceRoutes.GET("", controllers.GetInvoice)
	}

	captchaRoutes := router.Group("/captcha")
	{
		// Handle GET requests at /captcha
		captchaRoutes.GET("/:name", controllers.GetCaptcha)
		captchaRoutes.GET("/", controllers.ReloadCaptcha)
	}
}
