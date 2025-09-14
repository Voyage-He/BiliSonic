package bilibili

import (
	"log"
	"net/url"
	"testing"
)

// func createMockServer(handlerFunc http.HandlerFunc) *httptest.Server {
// 	return httptest.NewServer(handlerFunc)
// }

// func mockHandler(w http.ResponseWriter, r *http.Request) {
// 	switch {
// 	case strings.Contains(r.URL.Path, "/"):
// 		log.Println("cookie")
// 	case strings.Contains(r.URL.Path, "/search/type"):
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{
// 			"data": {
// 				"result": [
// 					{"title": "Video 1", "arcurl": "https://www.bilibili.com/video/1"},
// 					{"title": "Video 2", "arcurl": "https://www.bilibili.com/video/2"}
// 				]
// 			}
// 		}`))
// 	case strings.Contains(r.URL.Path, "/web-interface/view"):
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{
// 			"data": {
// 				"bvid": "BV1xK4y1u7gQ",
// 				"title": "Demo Video",
// 				"owner": {"name": "Author"},
// 				"pic": "https://example.com/thumb.jpg"
// 			}
// 		}`))
// 	case strings.Contains(r.URL.Path, "/player/playurl"):
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{
// 			"data": {
// 				"dash": {
// 					"audio": [
// 						{"baseUrl": "https://example.com/audio.mp3"}
// 					]
// 				}
// 			}
// 		}`))
// 	default:
// 		w.WriteHeader(http.StatusNotFound)
// 	}
// }

func TestGetCookie(t *testing.T) {
	client := NewBilibiliClient()
	host, _ := url.Parse("https://www.bilibili.com")
	cookie := client.Client.Jar.Cookies(host)

	log.Println(cookie[0])
}

func TestSearch(t *testing.T) {
	client := NewBilibiliClient()

	results, err := client.Search("祖娅纳惜")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	log.Println(results[0].Duration)
}

func TestGetVideoInfo(t *testing.T) {
	client := NewBilibiliClient()

	info, err := client.GetVideoInfo("14aYKzhEwG")
	if err != nil {
		t.Fatalf("GetVideoInfo failed: %v", err)
	}

	log.Println(info)
}

// func TestGetCovorArt(t *testing.T) {
// 	client := NewBilibiliClient()

// 	read, err := client.GetCoverArt("//i1.hdslb.com/bfs/archive/803060ec1cca5f294d2bb93afc4827f948721c69.jpg")

// 	if err != nil {
// 		t.Fatalf("GetVideoInfo failed: %v", err)
// 	}
// }

func TestGetCid(t *testing.T) {
	client := NewBilibiliClient()
	// t.Fatal("sdfds")
	s, err := client.getCid("14aYKzhEwG")
	if err != nil {
		log.Fatalln(err)
	}

	// if err != nil {
	// 	t.Fatalf("GetVideoInfo faile======d: %v", err)
	// }
	log.Println(s)
}

func TestGetAudioUrl(t *testing.T) {
	client := NewBilibiliClient()
	cid, _ := client.getCid("14aYKzhEwG")
	audioUrl, err := client.GetAudioUrl("14aYKzhEwG", cid)
	if err != nil {
		t.Fatalf("GetAudioUrl failed: %v", err)
	}

	log.Println(audioUrl)
}

func TestGetStream(t *testing.T) {
	client := NewBilibiliClient()
	stream, _, err := client.GetAudioStream("14aYKzhEwG")
	if err != nil {
		t.Fatalf("GetAudioUrl failed: %v", err)
	}
	defer stream.Close()

	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := stream.Read(buf)
		if n > 0 {
			log.Println(buf[:n])
		}
		if err != nil {
			log.Println("over!!!!")
			return
		}
	}
}
