package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowRegisterPage(c *gin.Context) {
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"register.html",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title": "eInvoice",
		},
	)
}
