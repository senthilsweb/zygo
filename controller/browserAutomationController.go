package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/senthilsweb/zygo/pkg/utils"
	"github.com/tidwall/gjson"

	"github.com/go-rod/rod"
)

func Export2PDF(c *gin.Context) {
	request_body := utils.GetStringFromGinRequestBody(c)
	webpage := gjson.Get(request_body, "message.webpage")
	filename := gjson.Get(request_body, "message.filename")
	page := rod.New().MustConnect().MustPage(webpage.String())
	response := page.MustWaitLoad().MustPDF()

	c.Writer.Header().Add("Content-type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename.String()+".pdf")
	c.Writer.Write(response)
	return
}

func Export2PNG(c *gin.Context) {
	request_body := utils.GetStringFromGinRequestBody(c)
	webpage := gjson.Get(request_body, "message.webpage")
	filename := gjson.Get(request_body, "message.filename")
	fullpage := gjson.Get(request_body, "message.fullpage")
	page := rod.New().MustConnect().MustPage(webpage.String())
	var response []byte
	if fullpage.Bool() {
		response = page.MustWaitLoad().MustScreenshotFullPage()
	} else {
		response = page.MustWaitLoad().MustScreenshot()
	}
	c.Writer.Header().Add("Content-type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename.String()+".png")
	c.Writer.Write(response)
	return
}
