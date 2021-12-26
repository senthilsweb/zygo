package controller

import (
	"context"

	"github.com/senthilsweb/zygo/pkg/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/slack"
	"github.com/tidwall/gjson"
)

func NotifySlack(c *gin.Context) {

	request_body := utils.GetStringFromGinRequestBody(c)

	token := utils.GetValElseSetEnvFallback(request_body, "SLACK_TOKEN")
	channel := utils.GetValElseSetEnvFallback(request_body, "SLACK_CHANNEL")
	log.Infof("token = [" + token + "]")
	log.Infof("channel = [" + channel + "]")
	notifier := notify.New()

	//Get the teams webhook url from the config
	//token := viper.GetString("notification.slack.token")
	//log.Infof("token = [" + token + "]")
	//Get the teams webhook url from the config
	//channel := viper.GetString("notification.slack.channel")

	// Provide your Slack OAuth Access Token
	slackService := slack.New(token)

	// Passing a Slack channel id as receiver for our messages.
	// Where to send our messages.
	slackService.AddReceivers(channel)

	// Tell our notifier to use the Slack service. You can repeat the above process
	// for as many services as you like and just tell the notifier to use them.
	notifier.UseServices(slackService)

	/* The commented code below is just another example how to read dynamic json using viper
	however we commented out for consistencies
	dyn_viper := viper.New()
	dyn_viper.SetConfigType("yaml")
	dyn_viper.ReadConfig(c.Request.Body)
	//subject := dyn_viper.GetString("subject")
	//body := dyn_viper.GetString("body")
	*/
	subject := gjson.Get(request_body, "message.subject")
	body := gjson.Get(request_body, "message.body")
	log.Info(subject)
	// Send a message
	err := notifier.Send(
		context.Background(),
		subject.String(),
		body.String(),
	)

	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"success": "false", "message": err})
		//c.AbortWithStatus(500)
		return
	}

	c.JSON(200, gin.H{"success": "true", "message": "Slack notification was successful"})
	return
}

func PostMessageInPrivateChannel(payload string) {

	request_body := payload

	token := utils.GetValElseSetEnvFallback(request_body, "SLACK_TOKEN")
	channel := utils.GetValElseSetEnvFallback(request_body, "SLACK_CHANNEL")
	notifier := notify.New()
	slackService := slack.New(token)
	slackService.AddReceivers(channel)
	// Tell our notifier to use the Slack service. You can repeat the above process
	// for as many services as you like and just tell the notifier to use them.
	notifier.UseServices(slackService)
	subject := gjson.Get(request_body, "message.subject")
	body := gjson.Get(request_body, "message.payload")

	// Send a message
	err := notifier.Send(
		context.Background(),
		subject.String(),
		body.String(),
	)

	if err != nil {
		log.Info("Slack notification failed")
		log.Fatal(err)
	}
	log.Info("Slack notification was successful")
	return
}
