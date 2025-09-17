package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"

	"example/subsonic/bilibili"

	"github.com/gin-gonic/gin"
)

// 根节点 <subsonic-response ...>
type SubsonicResponseXML struct {
	XMLName       xml.Name          `xml:"subsonic-response"`
	Status        string            `xml:"status,attr"`
	Version       string            `xml:"version,attr"`
	Xmlns         string            `xml:"xmlns,attr"`
	Type          string            `xml:"type,attr,omitempty"`
	ServerVersion string            `xml:"serverVersion,attr,omitempty"`
	OpenSubsonic  *bool             `xml:"openSubsonic,attr,omitempty"`
	Ping          *PingXML          `xml:"ping,omitempty"`
	SearchResult2 *SearchResultXML  `xml:"searchResult2,omitempty"`
	SearchResult3 *SearchResultXML  `xml:"searchResult3,omitempty"`
	Song          *SongXML          `xml:"song,omitempty"`
	Starred2      *SearchResultXML  `xml:"starred2,omitempty"`
	Error         *SubsonicErrorXML `xml:"error,omitempty"`
}

type SubsonicErrorXML struct {
	Code    int    `xml:"code,attr"`
	Message string `xml:"message,attr"`
}

type PingXML struct {
	// 空结构，ping 无子字段
}

type SearchResultXML struct {
	Artist []interface{} `xml:"artist,omitempty"`
	Album  []interface{} `xml:"album,omitempty"`
	Song   []SongXML     `xml:"song,omitempty"`
}

type SongXML struct {
	ID          string `xml:"id,attr"`
	IsDir       bool   `xml:"isDir,attr,omitempty"`
	Title       string `xml:"title,attr"`
	Artist      string `xml:"artist,attr,omitempty"`
	CoverArt    string `xml:"coverArt,attr,omitempty"`
	ContentType string `xml:"contentType,attr,omitempty"`
	Suffix      string `xml:"suffix,attr,omitempty"`
	Duration    int    `xml:"duration,attr,omitempty"`
	ArtistID    string `xml:"artistId,attr,omitempty"`
	Type        string `xml:"type,attr,omitempty"`
	IsVideo     bool   `xml:"isVideo,attr,omitempty"`
}

// 从 bilibili.BilibiliVideo 转成 SongXML
func SongFromXML(v *bilibili.BilibiliVideo) SongXML {
	return SongXML{
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

// 构造一个最顶层的“ok”响应
func createSubsonicOkResponseXML() SubsonicResponseXML {
	open := true
	return SubsonicResponseXML{
		Status:        "ok",
		Version:       VERSION,
		Xmlns:         "http://subsonic.org/restapi",
		Type:          "voyage",
		ServerVersion: SERVER_VERSION,
		OpenSubsonic:  &open,
	}
}

// 简单鉴权：u=voyage, p=141592
func authMiddlewareXML(c *gin.Context) {
	u := c.Query("u")
	p := c.Query("p")
	if u != "voyage" || p != "141592" {
		resp := SubsonicResponseXML{
			Status:  "failed",
			Version: VERSION,
			Xmlns:   "http://subsonic.org/restapi",
			Error: &SubsonicErrorXML{
				Code:    40,
				Message: "Bad username or password",
			},
		}
		c.XML(http.StatusUnauthorized, resp)
		c.Abort()
		return
	}
	c.Next()
}

// ping.view
func PingHandlerXML(c *gin.Context) {
	log.Println("ping invoke")
	resp := createSubsonicOkResponseXML()
	resp.Ping = &PingXML{}
	c.XML(http.StatusOK, resp)
}

// search2.view
func Search2HandlerXML(c *gin.Context) {
	log.Println("search2 invoke")
	resp := createSubsonicOkResponseXML()
	resp.SearchResult2 = &SearchResultXML{
		Artist: []interface{}{},
		Album:  []interface{}{},
		Song:   []SongXML{},
	}
	c.XML(http.StatusOK, resp)
}

// search3.view
func Search3HandlerXML(c *gin.Context) {
	log.Println("search3 invoke")
	cliAny, _ := c.Get("client")
	client := cliAny.(*bilibili.BilibiliClient)
	q := c.Query("query")

	videos, err := client.Search(q)
	if err != nil {
		log.Println("search error:", err)
		resp := SubsonicResponseXML{
			Status:  "failed",
			Version: VERSION,
			Xmlns:   "http://subsonic.org/restapi",
			Error: &SubsonicErrorXML{
				Code:    50,
				Message: err.Error(),
			},
		}
		c.XML(http.StatusInternalServerError, resp)
		return
	}

	songs := make([]SongXML, 0, len(videos))
	for _, v := range videos {
		songs = append(songs, SongFromXML(&v))
	}

	resp := createSubsonicOkResponseXML()
	resp.SearchResult3 = &SearchResultXML{Song: songs}
	c.XML(http.StatusOK, resp)
}

func GetSongXML(c *gin.Context) {
	log.Println("getSong invoke")
	cliAny, _ := c.Get("client")
	client := cliAny.(*bilibili.BilibiliClient)
	id := c.Query("id")

	video, err := client.GetVideoInfo(id)
	if err != nil {
		log.Println("search error:", err)
		resp := SubsonicResponseXML{
			Status:  "failed",
			Version: VERSION,
			Xmlns:   "http://subsonic.org/restapi",
			Error: &SubsonicErrorXML{
				Code:    50,
				Message: err.Error(),
			},
		}
		c.XML(http.StatusInternalServerError, resp)
		return
	}

	videoXML := SongFromXML(video)

	resp := createSubsonicOkResponseXML()
	resp.Song = &videoXML
	c.XML(http.StatusOK, resp)
}

// starred.view
func StarredHandlerXML(c *gin.Context) {
	log.Println("starred invoke")
	resp := createSubsonicOkResponseXML()
	resp.Starred2 = &SearchResultXML{
		Song: []SongXML{},
	}
	c.XML(http.StatusOK, resp)
}

// getCoverArt.view (GET)
func GetCoverArtHandlerXML(c *gin.Context) {
	log.Println("getCoverArt invoke")
	cliAny, _ := c.Get("client")
	client := cliAny.(*bilibili.BilibiliClient)
	id := c.Query("id")
	log.Println("id:" + id)
	if id == "al-" {
		return
	}

	file, err := client.GetCoverArt(id)
	if err != nil {
		log.Println("coverArt error:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	c.Header("Content-Type", "image/jpeg")
	c.Stream(func(w io.Writer) bool {
		buf := make([]byte, 32*1024)
		n, err := file.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			return true
		}
		return err == nil
	})
}

// getCoverArt.view (HEAD)
func HeadCoverArtHandlerXML(c *gin.Context) {
	log.Println("headCoverArt invoke")
	cliAny, _ := c.Get("client")
	client := cliAny.(*bilibili.BilibiliClient)
	id := c.Query("id")

	file, err := client.GetCoverArt(id)
	if err != nil {
		log.Println("coverArt head error:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	c.Status(http.StatusOK)
}

// stream.view
func StreamHandlerXML(c *gin.Context) {
	log.Println("stream invoke")
	cliAny, _ := c.Get("client")
	client := cliAny.(*bilibili.BilibiliClient)
	id := c.Query("id")

	file, _, err := client.GetAudioStream(id)
	if err != nil {
		log.Println("stream error:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	c.Header("Content-Type", "audio/mpeg")
	c.Stream(func(w io.Writer) bool {
		buf := make([]byte, 32*1024)
		n, err := file.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			return true
		}
		return err == nil
	})
}
