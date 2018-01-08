// handlers.invoice.go

package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"erpvietnam/ehoadon-website/settings"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func GetInvoice(c *gin.Context) {
	invoiceID := c.Query("s")
	CaptchaInput := c.Query("captchaInput")
	CaptchaID := c.Query("captchaID")

	if invoiceID == "" || CaptchaInput == "" || CaptchaID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Validate the captcha
	if !captcha.VerifyString(CaptchaID, CaptchaInput) {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"data": "CaptchaErr"})
		return
	}

	// Check if the invoice ID is valid
	FileName := fmt.Sprintf("%s%s.pdf", settings.Settings.InvoiceFilePath, invoiceID)

	fInfo, err := os.Stat(FileName)
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"data": "FileNotFound"})
			c.AbortWithStatus(http.StatusNotFound)
		}
	}

	sizePDF := fInfo.Size()
	buf := make([]byte, sizePDF)

	Openfile, err := os.Open(FileName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer Openfile.Close()
	// read file content into buffer
	fReader := bufio.NewReader(Openfile)
	fReader.Read(buf)

	// convert the buffer bytes to base64 string - use buf.Bytes() for new image
	pdfBase64 := base64.StdEncoding.EncodeToString(buf)

	c.JSON(http.StatusOK, pdfBase64)
	return
}
