package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("gin.log")

	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

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
	router.StaticFile("/index.html", "./index.html")

	router.Run() // listen and serve on 0.0.0.0:8080
}
