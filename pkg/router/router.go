package router

import (
	"io"
	"log"
	"os"

	"github.com/senthilsweb/zygo/controller"
	"github.com/senthilsweb/zygo/pkg/middleware"
	"github.com/senthilsweb/zygo/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Setup function
func Setup() *gin.Engine {
	r := gin.New()
	f, _ := os.Create(utils.AppExecutionPath() + "/" + os.Args[0] + ".log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	log.Println("Bootstrapping gin middlewares")
	// Middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.GinContextToContextMiddleware())
	log.Println("Setting up routes")
	r.GET("/api/ping", func(c *gin.Context) {
		pod_node_name := os.Getenv("NODE_NAME")

		if len(pod_node_name) == 0 {
			pod_node_name = "NIL"
		}

		c.JSON(200, gin.H{
			"message": "pong",
			"node":    pod_node_name,
		})
	})

	r.POST("/api/notify/slack", controller.NotifySlack)
	r.POST("/api/notify/mailgun", controller.NotifyMailgun)
	r.POST("/api/pdf/export", controller.Export2PDF)
	r.POST("/api/png/export", controller.Export2PNG)
	r.POST("/api/word/export", controller.Export2Word)
	r.POST("/api/redis/enqueue", controller.Enqueue)
	r.POST("/api/redis/publish", controller.Publish)
	r.GET("/api/redis/dequeue/:key", controller.Dequeue)
	r.GET("/api/ev/:key", controller.GetEnvironment)

	r.POST("/api/redis/hook/swissknife", controller.Swissknife)

	log.Println("Finished router setup")
	return r
}
