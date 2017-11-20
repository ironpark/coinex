package main

import (
	"github.com/ironpark/coinex/bucket"
	"time"
	"github.com/ironpark/coinex/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"io"
	//"strconv"
	"github.com/gin-contrib/sse"
	"github.com/ironpark/coinex/strategy"
	"fmt"
)

func main() {
	//load configs
	//init bucket
	buck := bucket.Instance()
	go buck.Run()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	s:=strategy.LoadStrategy("/Users/ironpark/GoglandProjects/strategyTest/main")
	s.Init()
	s.Init()
	fmt.Println(s.Info())
	fmt.Println(s.Info())

	s.KillProcess()
	cdb := db.Default()
	//APIS
	v1 := router.Group("/api/v1/")
	{
		//ohlc data
		v1.Handle(http.MethodGet,"/ticker/:market/:quote/:base/:res", func(c *gin.Context) {
			market := c.Param("market")
			quote := c.Param("quote")
			base := c.Param("base")
			res := c.Param("res")


			ohlc := cdb.GetOHLC(market,db.Pair{Quote:quote,Base:base},time.Now().Add(-time.Hour*24),time.Now(),res)
			c.JSON(http.StatusOK,ohlc)
		})

		//get supported markets
		v1.Handle(http.MethodGet,"/markets", func(c *gin.Context) {

		})

		//Server-Sent-Event for bucket status updates
		v1.Handle(http.MethodGet, "/sse/bucket", func(c *gin.Context) {

			listener := make(chan sse.Event)
			//func(Ex,Pair,Type string,ListenerID int64)

			sub := func(market string, pair db.Pair, status bucket.AssetStatus) {
				nType := "RealTime"
				if time.Now().Unix() - status.Last > 60*10 {
					nType = "BackFills"
				}

				listener <- sse.Event{
					Id:    "",
					Event: "message",
					Data: map[string]interface{}{
						"Exchange": market,
						"Base": pair.Base,
						"Quote": pair.Quote,
						"Stop":status.IsStop,
						"First": status.First,
						"Last":status.Last,
						"Type":nType,
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
		v1.POST("/bucket", func(c *gin.Context) {
			ex := c.PostForm("ex")
			base := c.PostForm("base")
			pair := c.PostForm("pair")
			//start := c.PostForm("start")
			//unix,_ := strconv.ParseInt(start,0,64)
			buck.AddToTrack(ex,db.Pair{Quote:base,Base:pair})
		})

		//get tracking assets
		v1.GET("/bucket", func(c *gin.Context) {
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
