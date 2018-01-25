package controllers

import (
	"bytes"
	"erpvietnam/ehoadon-website/models"
	. "erpvietnam/ehoadon-website/settings"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type RegisterForm struct {
	Code                   string `form:"Code" binding:"required"`
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

func CheckCompanyCode(c *gin.Context) {
	code := c.Query("Code")

	if code == "" {
		c.JSON(
			http.StatusBadRequest,
			"false",
		)
		return
	}

	var err error
	models.DB, err = sqlx.Connect(Settings.Database.DriverName, Settings.GetDbConn())
	if err != nil {
		log.Fatal(err)
	}
	defer models.DB.Close()

	client := models.Client{}
	err = client.GetByCode(code)

	if err == models.ErrClientNotFound {
		c.JSON(
			http.StatusOK,
			"true",
		)
		return
	} else if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		nil,
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

	//session := sessions.Default(c)
	//var lastRegister time.Time
	//v := session.Get("lastRegister")
	//if v == nil {
	//	lastRegister = time.Now()
	//} else {
	//	lastRegister = v.(time.Time)
	//}
	//session.Set("lastRegister", lastRegister)
	//session.Save()

	activeCode, _ := uuid.NewV4()
	client := models.Client{
		Code:                        strings.ToUpper(register.Code),
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
		mail := models.NewMail(client.Email, "Thông báo V/v đăng ký sử dụng hóa đơn điện tử eInvoice", "")
		err = mail.ParseTemplate("./templates/mailActive.html", client)
		if err != nil {
			ServerInternalError505(c)
			return
		}
		ok, _ := mail.SendEmail()

		if !ok {
			ServerInternalError505(c)
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
		if len(transInfo.ValidationErrors) > 0 {
			c.HTML(
				http.StatusOK,
				"register.html",
				gin.H{
					"title":            "eInvoice",
					"Register":         register,
					"Validation":       false,
					"ValidationErrors": transInfo.ValidationErrors,
					"CaptchaID":        captcha.New(),
				},
			)
		} else {
			ServerInternalError505(c)
		}
	}
}

func RegisterActive(c *gin.Context) {
	activeCode := c.Query("active_code")

	if activeCode == "" {
		activeCode = c.Param("active_code")
	}

	if activeCode == "" {
		PageNotFound404(c)
		return
	}

	var err error
	models.DB, err = sqlx.Connect(Settings.Database.DriverName, Settings.GetDbConn())
	if err != nil {
		log.Fatal(err)
	}
	defer models.DB.Close()

	client := models.Client{}
	transInfo := client.Active(activeCode)

	if !transInfo.ReturnStatus {
		notFound, _ := models.InArray(models.ErrClientActiveCodeNotFound, transInfo.ReturnError)
		if notFound {
			PageNotFound404(c)
			return
		}

		codeExpired, _ := models.InArray(models.ErrClientActiveCodeExpired, transInfo.ReturnError)
		if codeExpired {
			//TODO : Hien thi thong bao va gui lai mail
			return
		}
	}

	//active
}

func RegisterInitDB(c *gin.Context) {
	clientID := c.Query("client_id")

	if clientID == "" {
		clientID = c.Param("client_id")
	}

	if clientID == "" {
		c.String(
			http.StatusBadRequest,
			"",
		)
		return
	}
	var err error
	models.DB, err = sqlx.Connect(Settings.Database.DriverName, Settings.GetDbConn())
	if err != nil {
		log.Fatal(err)
	}
	defer models.DB.Close()

	id, _ := strconv.ParseInt(clientID, 10, 64)

	client := models.Client{}
	_ = client.Get(id)
	initDB := client.GetInitDB()

	t, err := template.ParseFiles("./templates/initdb.sql")
	if err != nil {
		c.String(
			http.StatusInternalServerError,
			"",
		)
		return
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, initDB); err != nil {
		c.String(
			http.StatusInternalServerError,
			"",
		)
		return
	}

	c.String(
		http.StatusOK,
		buf.String(),
	)
}
