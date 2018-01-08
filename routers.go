// routes.go

package main

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func initializeRoutes() {

	router.GET("/", showIndexPage)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	invoiceRoutes := router.Group("/invoices")
	{
		// Handle GET requests at /invoices
		invoiceRoutes.GET("", GetInvoice)
	}

	captchaRoutes := router.Group("/captcha")
	{
		// Handle GET requests at /captcha
		captchaRoutes.GET("/:name", GetCaptcha)
		captchaRoutes.GET("/", ReloadCaptcha)
	}
}

func showIndexPage(c *gin.Context) {
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title":             "eInvoice",
			"SearchCaptchaId":   captcha.New(),
			"RegisterCaptchaId": captcha.New(),
		},
	)
}
