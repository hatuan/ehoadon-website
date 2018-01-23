package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServerInternalError505(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"505.html",
		gin.H{
			"title": "eInvoice",
		},
	)

}
