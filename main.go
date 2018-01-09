package main

import (
	"io"
	"net/http"
	"os"

	"erpvietnam/ehoadon-website/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("gin.log")

	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := routers.GetRoute()

	router.Run() // listen and serve on 0.0.0.0:8080
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, templateName string) {

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data)
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data)
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}
