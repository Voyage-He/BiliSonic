package main

import (
	"io"
	"log"
	"net/http"

	"example/subsonic/bilibili"

	"github.com/gin-gonic/gin"
)

const VERSION = "1.16.1"
const SERVER_VERSION = "0.0.1"

type SubsonicResponse struct {
	Status        string       `json:"status"`
	Version       string       `json:"version"`
	Type          string       `json:"type"`
	ServerVersion string       `json:"serverVersion"`
	OpenSubsonic  bool         `json:"openSubsonic"`
	SearchResult2 SearchResult `json:"searchResult2"`
	SearchResult3 SearchResult `json:"searchResult3"`
	Starred2      Starred      `json:"starred2"`
}

type Response struct {
	SubsonicResponse SubsonicResponse `json:"subsonic-response"`
}

type SearchResult struct {
	Artist interface{}
	Album  interface{}
	Song   []Song `json:"song"`
}

type Starred = SearchResult

type Song struct {
	ID          string `json:"id"`
	IsDir       bool   `json:"isDir"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	CoverArt    string `json:"coverArt"`
	ContentType string `json:"contentType"`
	Suffix      string `json:"suffix"`
	Duration    int    `json:"duration"`
	ArtistID    string `json:"artistId"`
	Type        string `json:"type"`
	IsVideo     bool   `json:"isVideo"`
}

func SongFrom(v *bilibili.BilibiliVideo) Song {
	return Song{
		ID:          v.ID,
		IsDir:       false,
		Title:       v.Title,
		Artist:      v.Author,
		CoverArt:    v.Pic,
		ContentType: "audio/mpeg",
		Suffix:      "mp3",
		Duration:    v.Duration,
		ArtistID:    v.Author,
		Type:        "music",
		IsVideo:     false,
	}
}

func createSubsonicOkResponse() Response {
	return Response{
		SubsonicResponse: SubsonicResponse{
			Status:        "ok",
			Version:       VERSION,
			Type:          "voyage",
			ServerVersion: SERVER_VERSION,
			OpenSubsonic:  true,
		},
	}

}

func checkAuth(r *http.Request) bool {
	user := r.URL.Query().Get("u")
	pass := r.URL.Query().Get("p")
	// TODO: 实现更复杂的 token/md5 验证
	return user == "voyage" && pass == "141592"
}

func PingHandler(c *gin.Context) {
	log.Println("ping invoke")
	res := createSubsonicOkResponse()
	c.JSON(http.StatusOK, res)

}

func Search2Handler(c *gin.Context) {
	log.Println("search2API invoke")
	res := createSubsonicOkResponse()
	res.SubsonicResponse.SearchResult2 = SearchResult{
		Artist: []interface{}{},
		Album:  []interface{}{},
		Song:   []Song{},
	}
	c.JSON(http.StatusOK, res)
}

func Search3Handler(c *gin.Context) {
	log.Println("search3 invoke")
	client0, _ := c.Get("client")
	client, _ := client0.(*bilibili.BilibiliClient)
	q := c.Query("query")

	videos, _ := client.Search(q)
	songs := []Song{}
	for _, item := range videos {
		songs = append(songs, SongFrom(&item))
	}

	res := createSubsonicOkResponse()
	res.SubsonicResponse.SearchResult3 = SearchResult{
		Song: songs,
	}

	c.JSON(http.StatusOK, res)
}

func getCoverArtHandler(c *gin.Context) {
	log.Println("getcover invoke")
	client0, _ := c.Get("client")
	client, _ := client0.(*bilibili.BilibiliClient)
	id := c.Query("id")
	log.Println(id)

	file, err := client.GetCoverArt(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer file.Close()

	// 将文件内容写入响应
	c.Stream(func(w io.Writer) bool {
		buf := make([]byte, 32*1024) // 32KB buffer
		n, err := file.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			return true // 告诉 Gin 还有数据，下次继续调用
		}
		return err == nil
	})
}

func headCoverArtHandler(c *gin.Context) {
	log.Println("getcover invoke")
	client0, _ := c.Get("client")
	client, _ := client0.(*bilibili.BilibiliClient)
	id := c.Query("id")
	log.Println(id)

	file, err := client.GetCoverArt(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer file.Close()

	// 将文件内容写入响应
	c.Stream(func(w io.Writer) bool {
		buf := make([]byte, 32*1024) // 32KB buffer
		n, err := file.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			return true // 告诉 Gin 还有数据，下次继续调用
		}
		return err == nil
	})
}

func streamHandler(c *gin.Context) {
	log.Println("stream invoke")
	client0, _ := c.Get("client")
	client, _ := client0.(*bilibili.BilibiliClient)
	id := c.Query("id")
	log.Println(id)

	file, contentLength, err := client.GetAudioStream(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer file.Close()
	log.Println(contentLength)
	c.Header("Content-Type", "audio/mpeg")
	c.Header("Content-Length", contentLength)

	c.Stream(func(w io.Writer) bool {
		buf := make([]byte, 32*1024) // 32KB buffer
		n, err := file.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			return true // 告诉 Gin 还有数据，下次继续调用
		}
		return err == nil
	})
}

func starredHandler(c *gin.Context) {
	log.Println("starred invoke")

	res := createSubsonicOkResponse()
	res.SubsonicResponse.Starred2 = Starred{
		Song: []Song{},
	}

	c.JSON(http.StatusOK, res)
}
