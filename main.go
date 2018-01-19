package main

import (
	"io"
	"log"
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
	log.SetOutput(gin.DefaultWriter)

	router := routers.GetRoute()

	router.Run() // listen and serve on 0.0.0.0:8080
}
