package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type BilibiliClient struct {
	Client *http.Client // HTTP 客户端

}

// NewBilibiliClient 创建一个新的 BilibiliClient
func NewBilibiliClient() *BilibiliClient {
	jar, _ := cookiejar.New(nil)
	client := &BilibiliClient{Client: &http.Client{
		Jar: jar,
	}}

	req, _ := http.NewRequest("GET", "https://www.bilibili.com", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.56")
	res, _ := client.Client.Do(req)
	defer res.Body.Close()

	log.Println(res.Header.Get("set-cookie"))

	return client
}

// Search 通过关键词搜索音频
func (client *BilibiliClient) Search(keyword string) ([]BilibiliVideo, error) {
	queryURL := "http://api.bilibili.com/x/web-interface/search/type"
	queryParams := url.Values{}
	queryParams.Add("keyword", keyword)
	queryParams.Add("search_type", "video")

	req, _ := http.NewRequest("GET", queryURL+"?"+queryParams.Encode(), nil)

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	results := jsonResponse["data"].(map[string]interface{})["result"].([]interface{})
	return BilibiliVideoModelFromList(results), nil
}

// GetVideoInfo 获取视频信息
func (client *BilibiliClient) GetVideoInfo(bvid string) ([]string, error) {
	queryURL := "https://api.bilibili.com/x/web-interface/view"
	queryParams := url.Values{}
	queryParams.Add("bvid", bvid)

	req, _ := http.NewRequest("GET", queryURL+"?"+queryParams.Encode(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0") // 设置常见的浏览器 User-Agent

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	data := jsonResponse["data"].(map[string]interface{})
	return []string{
		data["bvid"].(string),  // BVID
		data["title"].(string), // 视频标题
		data["owner"].(map[string]interface{})["name"].(string), // 作者
		data["pic"].(string), // 视频封面
	}, nil
}

func (client *BilibiliClient) GetCoverArt(coverArt string) (io.ReadCloser, error) {
	if strings.HasPrefix(coverArt, "//") {
		queryURL := "http:" + coverArt

		req, err := http.NewRequest("GET", queryURL, nil)
		if err != nil {
			return nil, fmt.Errorf("req")
		}

		resp, err := client.Client.Do(req)
		if err != nil {
			return nil, err
		}

		return resp.Body, nil
	}
	log.Fatalln("stringsS")
	return nil, fmt.Errorf("fsdf")
}

func (client *BilibiliClient) HeadCoverArt(coverArt string) (io.ReadCloser, error) {
	if strings.HasPrefix(coverArt, "//") {
		queryURL := "http:" + coverArt

		req, err := http.NewRequest("HEAD", queryURL, nil)
		if err != nil {
			return nil, fmt.Errorf("req")
		}

		resp, err := client.Client.Do(req)
		if err != nil {
			return nil, err
		}

		return resp.Body, nil
	}
	log.Fatalln("stringsS")
	return nil, fmt.Errorf("fsdf")
}

// GetAudioUrl 获取音频 URL
func (client *BilibiliClient) GetAudioUrl(bvid string, cid int) (string, error) {
	queryURL := "http://api.bilibili.com/x/player/playurl"
	queryParams := url.Values{}
	queryParams.Add("bvid", bvid)
	queryParams.Add("cid", strconv.Itoa(cid))
	queryParams.Add("fnval", "16")

	req, _ := http.NewRequest("GET", queryURL+"?"+queryParams.Encode(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0") // 设置常见的浏览器 User-Agent

	resp, err := client.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	audioUrl := jsonResponse["data"].(map[string]interface{})["dash"].(map[string]interface{})["audio"].([]interface{})[0].(map[string]interface{})["baseUrl"].(string)
	return audioUrl, nil
}

func (client *BilibiliClient) getCid(id string) (int, error) {
	queryURL := "http://api.bilibili.com/x/player/pagelist"
	queryParams := url.Values{}
	queryParams.Add("bvid", id)

	req, _ := http.NewRequest("GET", queryURL+"?"+queryParams.Encode(), nil)

	resp, err := client.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return 0, fmt.Errorf("failed to parse response: %v", err)
	}

	cidf := jsonResponse["data"].([]interface{})[0].(map[string]interface{})["cid"].(float64)
	cid := int(cidf) // TODO float convert to string
	return cid, nil
}

func (client *BilibiliClient) GetAudioStream(id string) (io.ReadCloser, string, error) {
	cid, _ := client.getCid(id)
	audioUrl, _ := client.GetAudioUrl(id, cid)
	audioUrl_Url, _ := url.Parse(audioUrl)

	req, _ := http.NewRequest("GET", audioUrl, nil)
	req.Header.Add("Referer", "https://www.bilibili.com")
	req.Header.Add("Host", audioUrl_Url.Host)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.56")

	client_new := http.Client{}
	resp, err := client_new.Do(req)
	if err != nil {
		return nil, "", err
	}

	contentLength := resp.Header.Get("Content-Length")

	return (resp.Body), contentLength, nil
}

// BilibiliVideo 示例模型定义
type BilibiliVideo struct {
	ID     string
	Title  string
	AVID   int
	Author string
	MID    int
	Pic    string
}

// BilibiliVideoModelFromList 将 JSON 转为 BilibiliVideoModel 的切片
func BilibiliVideoModelFromList(data []interface{}) []BilibiliVideo {
	result := []BilibiliVideo{}
	for _, item := range data {
		video := item.(map[string]interface{})
		if video["bvid"].(string) == "" {
			continue
		}
		bvid := strings.Split(video["bvid"].(string), "BV")[1]
		result = append(result, BilibiliVideo{
			ID:     bvid,
			Title:  removeHTMLTags(video["title"].(string)),
			AVID:   int(video["aid"].(float64)),
			Author: video["author"].(string),
			MID:    int(video["mid"].(float64)),
			Pic:    video["pic"].(string),
		})
	}
	return result
}

func removeHTMLTags(input string) string {
	// Regular expression to match HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}
