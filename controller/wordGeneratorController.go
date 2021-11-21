package controller

import (
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/senthilsweb/zygo/pkg/cetak"
	"github.com/senthilsweb/zygo/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type TemplateData struct {
	Title   string
	Content string
}

const BufferSize = 100

func Export2Word(c *gin.Context) {
	request_body := utils.GetStringFromGinRequestBody(c)
	template_name := gjson.Get(request_body, "message.template_name")
	file_name := gjson.Get(request_body, "message.filename")
	payload := gjson.Get(request_body, "message.payload")
	log.Info(payload)
	ts_string := strconv.Itoa(int(time.Now().Unix()))
	if len(template_name.String()) == 0 {
		c.JSON(500, gin.H{"success": "false", "message": "Missing template name"})
		return
	}
	output_file_name := ""
	if len(file_name.String()) == 0 {
		output_file_name = utils.GetFileNameWithoutExt(template_name.String()) + "_" + ts_string
	}
	output_file_name = utils.GetFileNameWithoutExt(file_name.String()) + "_" + ts_string
	file_ext := utils.GetFileExt(file_name.String())
	file_full_name := output_file_name + file_ext

	//bytes, err = srv.contents.ReadFile(srv.indexFilePath())

	d, err := cetak.New("./templates/" + template_name.String())
	if err != nil {
		c.JSON(500, gin.H{"success": "false", "message": err})
		return
	}

	err = d.Generate(payload.Value(), "./temp/"+file_full_name)
	if err != nil {
		c.JSON(500, gin.H{"success": "false", "message": err})
		return
	}

	dat, err := ioutil.ReadFile("./temp/" + file_full_name)

	if err != nil {
		c.JSON(500, gin.H{"success": "false", "message": err})
		return
	}

	c.Writer.Header().Add("Content-type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+file_full_name)
	c.Writer.Write(dat)
	c.JSON(200, gin.H{"success": "true", "message": "File exported successfully"})
	return
}
