package main

import (
	"example/subsonic/bilibili"
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
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	CoverArt    string `json:"coverArt"`
	ContentType string `json:"contentType"`
	Suffix      string `json:"suffix"`
	ArtistID    string `json:"artistId"`
	Type        string `json:"type"`
	IsVideo     bool   `json:"isVideo"`
}

func SongFrom(v *bilibili.BilibiliVideo) Song {
	return Song{
		ID:          v.ID,
		Title:       v.Title,
		Artist:      v.Author,
		CoverArt:    v.Pic,
		ContentType: "audio/mpeg",
		Suffix:      "mp3",
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
