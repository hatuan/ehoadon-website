package controllers

import (
	"erpvietnam/ehoadon-website/models"
	. "erpvietnam/ehoadon-website/settings"
	"log"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
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

	var err error
	models.DB, err = sqlx.Connect(Settings.Database.DriverName, Settings.GetDbConn())
	if err != nil {
		log.Fatal(err)
	}
	defer models.DB.Close()

	activeCode, _ := uuid.NewV4()
	client := models.Client{
		Description:                 register.Description,
		IsActivated:                 false,
		ActivatedCode:               activeCode.String(),
		CultureID:                   "vi-VN",
		AmountDecimalPlaces:         3,
		AmountRoundingPrecision:     decimal.NewFromFloat(0.001),
		UnitAmountDecimalPlaces:     3,
		UnitAmountRoundingPrecision: decimal.NewFromFloat(0.001),
		CurrencyLCYId:               0,
		VatNumber:                   register.VatNumber,
		GroupUnitCode:               "",
		VatMethodCode:               "",
		ProvinceCode:                "",
		DistrictsCode:               "",
		Address:                     register.CompanyAddress,
		AddressTransition:           register.AddressTransition,
		Telephone:                   register.Telephone,
		Email:                       register.Email,
		Fax:                         register.Fax,
		Website:                     register.Website,
		RepresentativeName:          register.RepresentativeName,
		RepresentativePosition:      register.RepresentativeName,
		ContactName:                 register.ContactName,
		Mobile:                      register.Mobile,
		BankAccount:                 "",
		BankName:                    "",
		TaxAuthoritiesID:            new(int64),
		Version:                     1,
		RecCreatedByID:              0,
		RecCreated:                  &models.Timestamp{time.Now()},
		RecModifiedByID:             0,
		RecModified:                 &models.Timestamp{time.Now()},
	}

	transInfo := client.Update()
	if transInfo.ReturnStatus {
		mail := models.NewMail([]string{client.Email, "hatuan05@gmail.com"}, "Thông báo V/v đăng ký sử dụng hóa đơn điện tử eInvoice", "")
		err = mail.ParseTemplate("./templates/mailActive.html", client)
		if err != nil {
			c.HTML(
				http.StatusInternalServerError,
				"register.html",
				gin.H{
					"title":     "eInvoice",
					"Register":  register,
					"CaptchaID": captcha.New(),
				},
			)
			return
		}
		ok, _ := mail.SendEmail()

		if !ok {
			c.HTML(
				http.StatusInternalServerError,
				"register.html",
				gin.H{
					"title":     "eInvoice",
					"Register":  register,
					"CaptchaID": captcha.New(),
				},
			)
		} else {
			c.HTML(
				http.StatusOK,
				"register.html",
				gin.H{
					"title":         "eInvoice",
					"Register":      register,
					"ShowActiveMsg": true,
					"CaptchaID":     captcha.New(),
				},
			)
		}
	} else {
		c.HTML(
			http.StatusInternalServerError,
			"register.html",
			gin.H{
				"title":     "eInvoice",
				"Register":  register,
				"CaptchaID": captcha.New(),
			},
		)
	}
}
