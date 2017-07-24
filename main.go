package main

import (
	"github.com/ironpark/coinex/bucket"
	"time"
	"github.com/ironpark/coinex/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"io"
	"runtime"
	"path"
	"strconv"
	"github.com/gin-contrib/sse"
)

func main() {
	//load configs
	//init bucket
	buck := bucket.Instance()
	go buck.Run()
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

			sub := func(market string, pair db.Pair, status bucket.AssetStatus) {
				listener <- sse.Event{
					Id:    "124",
					Event: "message",
					Data: map[string]interface{}{
						"Exchange": market,
						"Base": pair.Base,
						"Pair": pair.Quote,
						//"Stop":item.Stop,
						"First": status.First,
						"Last":status.Last,
						//"Type":Type,
					},
				}
			}
			buck.SubStatusUpdate(sub)
			c.Stream(func(w io.Writer) bool {
				//TODO First Data
				c.SSEvent("message", <-listener)
				return true
			})
			buck.UnSubStatusUpdate(sub)
		})
		//Server-Sent-Event for ticker data
		//v1.Handle(http.MethodGet, "/sse/ticker/:ex/:pair/:res", func(c *gin.Context) {
		//	ex := c.Param("name")
		//	pair := c.Param("pair")
		//	//res := c.Param("res")
		//
		//	listener := make(chan sse.Event)
		//	id := buck.AddEventListener(ex, func(Ex, Pair, Type string, ListenerID int64) {
		//		listener <- sse.Event{
		//			Id:    "124",
		//			Event: "message",
		//			Data: map[string]interface{}{
		//				Ex:   ex,
		//				Type: Type,
		//				Pair: Pair,
		//			},
		//		}
		//	},pair)
		//
		//	c.Stream(func(w io.Writer) bool {
		//		//TODO First Data
		//		c.SSEvent("message", <-listener)
		//		buck.RemoveEventListener(ex,id)
		//		return true
		//	})
		//})


		////add asset tracking
		v1.POST("/bucket/assets", func(c *gin.Context) {
			ex := c.PostForm("ex")
			base := c.PostForm("base")
			pair := c.PostForm("pair")
			start := c.PostForm("start")
			unix,_ := strconv.ParseInt(start,0,64)
			buck.AddToTrack(ex,db.Pair{Quote:base,Base:pair},time.Unix(unix,0))
		})

		//get tracking assets
		v1.GET("/bucket/assets", func(c *gin.Context) {
			c.JSON(http.StatusOK,buck.Status())
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
