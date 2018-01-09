package controllers

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func ShowIndexPage(c *gin.Context) {
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
