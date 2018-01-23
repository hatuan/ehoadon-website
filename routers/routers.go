// routers.go

package routers

import (
	"erpvietnam/ehoadon-website/controllers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()

	store := sessions.NewCookieStore([]byte("b2344aed-8ec3-41dc-964b-4da318a7475f"))
	router.Use(sessions.Sessions("ehoadon", store))

	router.LoadHTMLGlob("templates/*")
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

	// Initialize the routes
	initializeRoutes()
}

func GetRoute() *gin.Engine {
	return router
}

func initializeRoutes() {

	router.NoRoute(controllers.PageNotFound404)

	router.GET("/", controllers.ShowIndexPage)
	router.GET("/index.html", controllers.ShowIndexPage)
	router.GET("/index", controllers.ShowIndexPage)

	router.GET("/register", controllers.ShowRegisterPage)
	router.GET("/register.html", controllers.ShowRegisterPage)
	router.POST("/register", controllers.RegisterNewCompany)

	router.GET("/active", controllers.RegisterActive)
	router.GET("/active/:active_code", controllers.RegisterActive)

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
