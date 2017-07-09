package main

import (
	"github.com/ironpark/coinex/bucket"
	"time"
	"github.com/ironpark/coinex/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"io"
	"github.com/manucorporat/sse"
	"runtime"
	"path"
	"strconv"
	"log"
)

func main() {
	coindb,err := db.Default()
	if err != nil{
		log.Fatal(err)
	}
	//load configs
	localConfig := Config()
	//init bucket
	buck := bucket.NewBucket()
	for  _,asset:= range localConfig.Bucket.Assets {

		first, _ := coindb.FirstTradeHistory(asset.Base, asset.Pair)
		last, _ := coindb.LastTradeHistory(asset.Base, asset.Pair)
		log.Println(asset,first,last)
		buck.Add(&bucket.Target{
			Stop:     false,
			Exchange: asset.Ex,
			Base:     asset.Base,
			Pair:     asset.Pair,
			First:    first,
			Last:     last,
			Start:    last,
			End:      time.Now().UTC(),
		})
	}
	buck.Run()

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
			//if c.Request.URL.RawPath == "/"{
			//	c.Redirect(http.StatusPermanentRedirect,"/#/home")
			//}
			c.HTML(http.StatusOK, "index.html", nil)
		})
	}
	//APIS
	v1 := router.Group("/api/v1/")
	{
		//Server-Sent-Event for bucket status updates
		v1.Handle(http.MethodGet, "/sse/bucket", func(c *gin.Context) {
			listener := make(chan sse.Event)
			id := buck.AddGlobalEventListener(func(Ex, Type, Pair string, id int64) {
				fmt.Println(Ex, Type, Pair, id)
				listener <- sse.Event{
					Id:    "124",
					Event: "message",
					Data: map[string]interface{}{
						Ex:   Ex,
						Type: Type,
						Pair: Pair,
					},
				}
			})
			c.Stream(func(w io.Writer) bool {
				//TODO First Data
				c.SSEvent("message", <-listener)
				buck.RemoveGlobalEventListener(id)
				return true
			})
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

		//get bucket status
		v1.GET("/bucket/status", func(c *gin.Context) {
			c.JSON(http.StatusOK,buck.Assets)
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

		//get tracked assets
		v1.GET("/bucket/assets", func(c *gin.Context) {
			c.JSON(http.StatusOK,localConfig.Bucket.Assets)
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
