// handlers.invoice.go

package controllers

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

// GetCaptcha router to get the capche image and audio.
func GetCaptcha(c *gin.Context) {
	name := c.Param("name")
	var captchatype string
	var captchaid string
	if len(name) <= 4 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	captchatype = name[len(name)-3:]
	captchaid = name[:len(name)-4]
	if captchatype == "png" {
		err := captcha.WriteImage(c.Writer, captchaid, captcha.StdWidth, captcha.StdHeight)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

// ReloadCaptcha router to handle the post request of sign up page and get a new capche image and audio.
func ReloadCaptcha(c *gin.Context) {
	update := c.Param("update")
	if update == "true" {
		c.JSON(http.StatusOK, gin.H{
			"status":    "reloadCaptcha",
			"captchaid": captcha.New(),
		})
	}
}
