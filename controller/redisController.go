package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/senthilsweb/zygo/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func Enqueue(c *gin.Context) {

	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	request_body := utils.GetStringFromGinRequestBody(c)
	//redis_host := utils.GetValElseSetEnvFallback(request_body, "MAILGUN_DOMAIN")
	redis_uri := utils.GetValElseSetEnvFallback(request_body, "REDIS_URI")

	identity := gjson.Get(request_body, "message.identity")
	kv_key := gjson.Get(request_body, "message.kv_key").String()

	kv_body := gjson.Get(request_body, "message.kv_body").String()

	opt, _ := redis.ParseURL(redis_uri)
	client := redis.NewClient(opt)

	if identity.Bool() {
		kv_key = kv_key + ":" + uuid
	}

	// Publish a generated user to the new_users channel
	ctx := context.Background()
	log.Info("kv_key=" + kv_key)
	log.Info("kv_value=" + kv_body)

	client.Set(ctx, kv_key, kv_body, 0)
	c.JSON(200, gin.H{"success": "true", "message": "Successfully Enqueued", "key": kv_key})
	return
}

func Dequeue(c *gin.Context) {
	log.Info("Inside Dequeue")
	request_body := utils.GetStringFromGinRequestBody(c)
	kv_key := c.Param("key")
	log.Info("kv_key = " + kv_key)
	redis_uri := utils.GetValElseSetEnvFallback(request_body, "REDIS_URI")

	opt, _ := redis.ParseURL(redis_uri)
	client := redis.NewClient(opt)

	ctx := context.Background()

	val := client.Get(ctx, kv_key).Val()
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(val), &jsonMap)
	c.JSON(200, gin.H{"success": "true", "message": "Successfully Dequeued", "data": jsonMap})
	return
}

func List(c *gin.Context) {
	log.Info("List")
	request_body := utils.GetStringFromGinRequestBody(c)
	page := gjson.Get(request_body, "page").Int()
	key_prefix := gjson.Get(request_body, "key_prefix").String()
	redis_uri := utils.GetValElseSetEnvFallback(request_body, "REDIS_URI")

	if page <= 0 {
		c.JSON(500, gin.H{"success": "false", "message": "invalid page"})
		return
	}

	const pageSize int64 = 20
	from, to := (int64(page)-1)*pageSize, int64(page)*pageSize-1

	opt, _ := redis.ParseURL(redis_uri)
	client := redis.NewClient(opt)
	ctx := context.Background()

	ids, err := client.LRange(ctx, key_prefix, from, to).Result()
	if err != nil {
		c.JSON(500, gin.H{"success": "false", "message": "Failure", "exception": err})
	}

	var cmds []*redis.SliceCmd
	_, err = client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, keyId := range ids {
			cmds = append(cmds, pipe.HMGet(ctx, fmt.Sprintf(key_prefix+":%s", keyId)))
		}
		return nil
	})
	if err != nil {
		c.JSON(500, gin.H{"success": "false", "message": "Failure", "exception": err})
	}

	var jsonMap map[string]string
	var res []interface{}

	for _, c := range cmds {
		id := fmt.Sprint(c.Val())

		val := client.Get(ctx, id).Val()

		json.Unmarshal([]byte(val), &jsonMap)

		res = append(res, jsonMap)
	}

	c.JSON(200, gin.H{"success": "true", "message": "Successfully Dequeued", "data": res})
	return
}

func Swissknife(c *gin.Context) {
	request_body := utils.GetStringFromGinRequestBody(c)
	redis_uri := utils.GetValElseSetEnvFallback(request_body, "REDIS_URI")
	opt, _ := redis.ParseURL(redis_uri)
	client := redis.NewClient(opt)
	// Publish a generated user to the new_users channel
	ctx := context.Background()
	client.Set(ctx, "last_message", request_body, 0)

	number := gjson.Get(request_body, "number").String()
	form_name := gjson.Get(request_body, "form_name").String()
	kv_key := ""
	if len(number) > 0 && len(form_name) > 0 {
		kv_key = form_name + ":" + number
	}
	log.Info("kv_key=" + kv_key)
	client.Set(ctx, form_name+":"+number, request_body, 0)
	c.JSON(200, gin.H{"success": "true", "message": "Webhook payload was successfully Enqueued", "key": form_name + ":" + number})
	return
}

func GetEnvironment(c *gin.Context) {
	key := c.Param("key")
	val := os.Getenv(key)

	c.JSON(200, gin.H{"success": "true", "message": "Environment variable get attempt was successful", "key": key, "value": val})
	return
}

func Publish(c *gin.Context) {
	log.Info("Inside Publish")
	request_body := utils.GetStringFromGinRequestBody(c)
	channel := gjson.Get(request_body, "message.channel").String()
	//payload := gjson.Get(request_body, "message.payload").String()

	redis_uri := utils.GetValElseSetEnvFallback(request_body, "REDIS_URI")

	opt, _ := redis.ParseURL(redis_uri)
	client := redis.NewClient(opt)
	// Publish a generated user to the new_users channel
	ctx := context.Background()
	client.Set(ctx, "last_publish_message", request_body, 0)

	err := client.Publish(ctx, channel, request_body).Err()

	if err != nil {
		c.JSON(500, gin.H{"success": "false", "message": "Message publish failed.", "exception": err})
		return
	}

	log.Info("Message body " + request_body)

	c.JSON(200, gin.H{"success": "true", "message": "Message has been published successfully"})
	return
}

//This method is to receive the message sent by netlify to redis
func SubscribeAndReceiveMessage() {
	log.Info("Inside SubscribeAndReceiveMessage")
	redis_uri := os.Getenv("REDIS_URI")
	log.Info("redis_uri = " + redis_uri)
	opt, _ := redis.ParseURL(redis_uri)
	client := redis.NewClient(opt)
	ctx := context.Background()
	topic := client.Subscribe(ctx, "zynomi-website-lead")
	defer topic.Close()

	for {
		msg, err := topic.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
			log.Error(err)
		}
		log.Info("Message Received" + msg.Payload)
		log.Info("Posting message to slack")
		PostMessageInPrivateChannel(msg.Payload)
		log.Info("Posted message to slack")
		client.Set(ctx, "last_received_message", msg.Payload, 0)
	}
}
