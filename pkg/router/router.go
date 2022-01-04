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

	// hardcoded response for debugging.
	r.GET("/api/contacts", func(c *gin.Context) {
		jsonData := []byte(`[{"id":1,"first_name":"Alleen","last_name":"D'Alesco","email":"adalesco0@techcrunch.com","gender":"Polygender","address":"61 Crownhardt Lane"},{"id":2,"first_name":"Cassie","last_name":"De Angelo","email":"cdeangelo1@godaddy.com","gender":"Genderfluid","address":"59 West Road"},{"id":3,"first_name":"Mahala","last_name":"Scarlet","email":"mscarlet2@techcrunch.com","gender":"Bigender","address":"08 Victoria Circle"},{"id":4,"first_name":"Lesya","last_name":"Mantha","email":"lmantha3@wunderground.com","gender":"Agender","address":"832 Kedzie Place"},{"id":5,"first_name":"Minna","last_name":"Klemencic","email":"mklemencic4@wired.com","gender":"Male","address":"15 Amoth Court"},{"id":6,"first_name":"Libbey","last_name":"Du Barry","email":"ldubarry5@ebay.com","gender":"Genderqueer","address":"3 Pleasure Court"},{"id":7,"first_name":"Yves","last_name":"Soots","email":"ysoots6@jimdo.com","gender":"Agender","address":"23574 Macpherson Avenue"},{"id":8,"first_name":"Elnar","last_name":"Burne","email":"eburne7@blogs.com","gender":"Bigender","address":"0 Schmedeman Road"},{"id":9,"first_name":"Fiona","last_name":"Duell","email":"fduell8@delicious.com","gender":"Male","address":"49 Messerschmidt Circle"},{"id":10,"first_name":"Fee","last_name":"Esmonde","email":"fesmonde9@flickr.com","gender":"Male","address":"045 Paget Hill"},{"id":11,"first_name":"Kaine","last_name":"Elcott","email":"kelcotta@ibm.com","gender":"Genderfluid","address":"20671 Westridge Trail"},{"id":12,"first_name":"Leslie","last_name":"Luisetti","email":"lluisettib@icq.com","gender":"Male","address":"32772 Butterfield Road"},{"id":13,"first_name":"Mauricio","last_name":"Sappell","email":"msappellc@sphinn.com","gender":"Male","address":"5 Kenwood Street"},{"id":14,"first_name":"Alexandro","last_name":"Blunn","email":"ablunnd@ehow.com","gender":"Bigender","address":"651 Namekagon Trail"},{"id":15,"first_name":"Gus","last_name":"Alfonso","email":"galfonsoe@washingtonpost.com","gender":"Non-binary","address":"2062 Dwight Alley"},{"id":16,"first_name":"Derry","last_name":"Coppo","email":"dcoppof@free.fr","gender":"Non-binary","address":"49 Birchwood Park"},{"id":17,"first_name":"Loni","last_name":"Verbeek","email":"lverbeekg@themeforest.net","gender":"Polygender","address":"2572 Autumn Leaf Plaza"},{"id":18,"first_name":"Lydon","last_name":"Eglaise","email":"leglaiseh@cnet.com","gender":"Non-binary","address":"98995 Scofield Hill"},{"id":19,"first_name":"Claire","last_name":"McGucken","email":"cmcguckeni@blinklist.com","gender":"Female","address":"4757 Kipling Alley"},{"id":20,"first_name":"Reggie","last_name":"Hardstaff","email":"rhardstaffj@wp.com","gender":"Non-binary","address":"7 Haas Place"},{"id":21,"first_name":"Binni","last_name":"Swoffer","email":"bswofferk@1und1.de","gender":"Genderfluid","address":"43 Northridge Alley"},{"id":22,"first_name":"Sal","last_name":"Gillibrand","email":"sgillibrandl@narod.ru","gender":"Agender","address":"1439 Arkansas Trail"},{"id":23,"first_name":"Lois","last_name":"Torns","email":"ltornsm@xrea.com","gender":"Genderfluid","address":"84 Ilene Lane"},{"id":24,"first_name":"Cynthea","last_name":"Hellyer","email":"chellyern@unc.edu","gender":"Male","address":"0153 Springview Hill"},{"id":25,"first_name":"Charlene","last_name":"Arnoldi","email":"carnoldio@goo.ne.jp","gender":"Bigender","address":"42786 International Park"}]`)
		c.Data(200, "application/json", jsonData)
	})

	r.POST("/api/notify/slack", controller.NotifySlack)
	r.POST("/api/notify/mailgun", controller.NotifyMailgun)
	r.POST("/api/pdf/export", controller.Export2PDF)
	r.POST("/api/png/export", controller.Export2PNG)
	r.POST("/api/word/export", controller.Export2Word)
	r.POST("/api/redis/enqueue", controller.Enqueue)
	r.POST("/api/redis/publish", controller.Publish)
	r.GET("/api/redis/dequeue/:key", controller.Dequeue)
	r.POST("/api/redis/list", controller.List)
	r.GET("/api/ev/:key", controller.GetEnvironment)

	r.POST("/api/redis/hook/swissknife", controller.Swissknife)

	log.Println("Finished router setup")
	return r
}
