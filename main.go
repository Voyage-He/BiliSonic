package main

import (
	"log"

	"example/subsonic/bilibili"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	client := bilibili.NewBilibiliClient()
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("client", client)
		c.Next()
	})

	router.GET("/", PingHandler)

	router.HEAD("/", PingHandler)
	router.GET("/rest/ping.view", PingHandler)
	router.GET("/rest/search2.view", Search2Handler)
	router.GET("/rest/search3.view", Search3Handler)
	router.GET("/rest/getCoverArt.view", getCoverArtHandler)
	router.HEAD("/rest/getCoverArt.view", headCoverArtHandler)
	router.GET("/rest/stream.view", streamHandler)
	router.GET("rest/scrobble.view", PingHandler)
	router.GET("rest/getStarred2.view", starredHandler)

	router.GET("/rest/ping", PingHandlerXML)
	router.GET("/rest/search2", Search2HandlerXML)
	router.GET("/rest/search3", Search3HandlerXML)
	// router.GET("/rest/getCoverArt", GetCoverArtHandlerXML)
	router.GET("/rest/getSong", GetSongXML)
	// router.HEAD("/rest/getCoverArt", headCoverArtHandler)
	router.GET("/rest/stream", StreamHandlerXML)
	router.GET("rest/scrobble", PingHandlerXML)
	// router.GET("rest/getStarred2", starredHandlerXML)

	log.Println("OpenSubsonic proxy running at :8080")
	router.Run()
}
