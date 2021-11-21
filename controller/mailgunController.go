package controller

import (
	"context"
	"os"
	"time"

	"github.com/senthilsweb/zygo/pkg/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	mailgun "github.com/mailgun/mailgun-go/v4"
	"github.com/tidwall/gjson"
)

func NotifyMailgun(c *gin.Context) {

	pod_node_name := os.Getenv("NODE_NAME")

	if len(pod_node_name) == 0 {
		pod_node_name = "NIL"
	}
	log.Info("NODE_NAME = [" + pod_node_name + "]")

	request_body := utils.GetStringFromGinRequestBody(c)

	mailgun_domain := utils.GetValElseSetEnvFallback(request_body, "MAILGUN_DOMAIN")
	mailgun_key := utils.GetValElseSetEnvFallback(request_body, "MAILGUN_KEY")
	sender := utils.GetValElseSetEnvFallback(request_body, "EMAIL_SENDER")

	subject := gjson.Get(request_body, "message.subject")
	body := gjson.Get(request_body, "message.body")
	mailgun_email_template := gjson.Get(request_body, "message.template")
	mailgun_email_payload := gjson.Get(request_body, "message.payload")
	recipient := gjson.Get(request_body, "message.recipient")

	log.Info("mailgun_email_payload = " + mailgun_email_payload.String())
	log.Info("Body = " + body.String())

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(mailgun_domain, mailgun_key)

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject.String(), "", recipient.String())
	message.SetTemplate(mailgun_email_template.String())
	err := message.AddTemplateVariable("passwordResetLink", mailgun_email_payload.String())
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"success": "false", "message": err})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"success": "false", "message": err})
		return
	}
	log.Info(id)
	c.JSON(200, gin.H{"success": "true", "message": "Email has been sent successfully", "m": resp, "node": pod_node_name})
	return

}
