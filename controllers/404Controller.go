package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PageNotFound404(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"404.html",
		gin.H{
			"title": "eInvoice",
		},
	)
}
