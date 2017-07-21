package main

import (
	"github.com/ironpark/coinex/bucket"
	"time"
	"github.com/ironpark/coinex/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"io"
	"runtime"
	"path"
	"strconv"
	"log"
	"github.com/gin-contrib/sse"
)

func main() {
	//load configs
	localConfig := Config()
	//init bucket
	buck := bucket.GetInstance()
	ui := false
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	//init web ui
	if ui {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			panic("No caller information")
		}

		resource := path.Dir(filename) + "/front/dist"
		router.LoadHTMLGlob(resource + "/*.html")
		router.Static("/static", resource+"/static")
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})
	}

	//APIS
	v1 := router.Group("/api/v1/")
	{
		//Server-Sent-Event for bucket status updates
		v1.Handle(http.MethodGet, "/sse/bucket", func(c *gin.Context) {

			listener := make(chan sse.Event)
			//func(Ex,Pair,Type string,ListenerID int64)
			buck.Subscribe(bucket.TOPIC_UPDATE, func() {

			})
			id := buck.AddGlobalEventListener(func(Ex, Pair, Type string, id int64) {
				fmt.Println(Ex, Type, Pair, id)
				for _,item := range buck.GetStatus(){
					if item.Exchange == Ex{
						if item.Base + "_" + item.Pair == Pair {
							listener <- sse.Event{
								Id:    "124",
								Event: "message",
								Data: map[string]interface{}{
									"Stop":item.Stop,
									"Exchange": item.Exchange,
									"Base": item.Base,
									"Pair": item.Pair,
									"First": item.First,
									"Last":item.Last,
									"Type":Type,
								},
							}
						}
					}
				}
			})

			c.Stream(func(w io.Writer) bool {
				//TODO First Data
				c.SSEvent("message", <-listener)
				return true
			})
			buck.RemoveGlobalEventListener(id)

		})
		//Server-Sent-Event for ticker data
		v1.Handle(http.MethodGet, "/sse/ticker/:ex/:pair/:res", func(c *gin.Context) {
			ex := c.Param("name")
			pair := c.Param("pair")
			//res := c.Param("res")

			listener := make(chan sse.Event)
			id := buck.AddEventListener(ex, func(Ex, Pair, Type string, ListenerID int64) {
				listener <- sse.Event{
					Id:    "124",
					Event: "message",
					Data: map[string]interface{}{
						Ex:   ex,
						Type: Type,
						Pair: Pair,
					},
				}
			},pair)

			c.Stream(func(w io.Writer) bool {
				//TODO First Data
				c.SSEvent("message", <-listener)
				buck.RemoveEventListener(ex,id)
				return true
			})
		})


		//add asset tracking
		v1.POST("/bucket/assets", func(c *gin.Context) {
			ex := c.PostForm("ex")
			base := c.PostForm("base")
			pair := c.PostForm("pair")
			start := c.PostForm("start")
			unix,_ := strconv.ParseInt(start,0,64)
			localConfig.AddTarget(ex,base,pair,time.Unix(unix,0))

			buck.Add(&bucket.Target{
				Stop:     false,
				Exchange: ex,
				Base:     base,
				Pair:     pair,
				First:    time.Unix(unix,0),
				Last:     time.Now(),
				Start:    time.Unix(unix,0),
				End:      time.Now(),
			})
		})

		//get tracking assets
		v1.GET("/bucket/assets", func(c *gin.Context) {
			c.JSON(http.StatusOK,buck.GetStatus())
		})

		//update tracked assets setting
		v1.PUT("/bucket/assets/:id", func(c *gin.Context) {

		})

		//remove tracked assets
		v1.DELETE("/bucket/assets/:id", func(c *gin.Context) {

		})
	}
	router.Run()
}
