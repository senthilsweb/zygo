package controller

import (
	"context"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/slack"
)

func NotifySlack(c *gin.Context) {

	notifier := notify.New()

	//Get the teams webhook url from the config
	token := viper.GetString("notification.slack.token")
	log.Infof("token = [" + token + "]")
	//Get the teams webhook url from the config
	channel := viper.GetString("notification.slack.channel")

	// Provide your Slack OAuth Access Token
	slackService := slack.New(token)

	// Passing a Slack channel id as receiver for our messages.
	// Where to send our messages.
	slackService.AddReceivers(channel)

	// Tell our notifier to use the Slack service. You can repeat the above process
	// for as many services as you like and just tell the notifier to use them.
	notifier.UseServices(slackService)

	dyn_viper := viper.New()
	dyn_viper.SetConfigType("yaml")
	dyn_viper.ReadConfig(c.Request.Body)

	subject := dyn_viper.GetString("subject")
	body := dyn_viper.GetString("body")
	log.Info(subject)
	// Send a message
	err := notifier.Send(
		context.Background(),
		subject,
		body,
	)

	if err != nil {
		log.Fatal(err)
		c.AbortWithStatus(500)
		return
	}

	c.JSON(200, gin.H{"success": "true", "message": "Notification sent successful"})
	return
}
