package controllers

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

type RegisterForm struct {
	Description            string `form:"Description" binding:"required"`
	VatNumber              string `form:"VatNumber" binding:"required"`
	CompanyAddress         string `form:"CompanyAddress" binding:"required"`
	AddressTransition      string `form:"AddressTransition" binding:"required"`
	Telephone              string `form:"Telephone" binding:"required"`
	Fax                    string `form:"Fax" binding:"required"`
	Email                  string `form:"Email" binding:"required"`
	Website                string `form:"Website" binding:"required"`
	RepresentativeName     string `form:"RepresentativeName" binding:"required"`
	RepresentativePosition string `form:"RepresentativePosition" binding:"required"`
	ContactName            string `form:"ContactName" binding:"required"`
	Mobile                 string `form:"Mobile" binding:"required"`
	CaptchaInput           string `form:"CaptchaInput" binding:"required"`
	CaptchaID              string `form:"CaptchaID" binding:"required"`
}

func ShowRegisterPage(c *gin.Context) {

	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"register.html",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title":     "eInvoice",
			"Register":  RegisterForm{},
			"CaptchaID": captcha.New(),
		},
	)
}

func RegisterNewCompany(c *gin.Context) {

	register := RegisterForm{}
	c.Bind(&register)

	if register.CaptchaInput == "" || register.CaptchaID == "" {
		//redisplay register
		c.HTML(
			http.StatusBadRequest,
			"register.html",
			gin.H{
				"title":        "eInvoice",
				"CaptchaError": true,
				"Register":     register,
				"CaptchaID":    captcha.New(),
			},
		)
		return
	}

	// Validate the captcha
	if !captcha.VerifyString(register.CaptchaID, register.CaptchaInput) {
		c.HTML(
			http.StatusOK,
			"register.html",
			gin.H{
				"title":        "eInvoice",
				"CaptchaError": true,
				"Register":     register,
				"CaptchaID":    captcha.New(),
			},
		)
		return
	}

	c.HTML(
		http.StatusOK,
		"register.html",
		gin.H{
			"title":     "eInvoice",
			"Register":  register,
			"CaptchaID": captcha.New(),
		},
	)
}
