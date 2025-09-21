package main

import (
	"encoding/json"
	"os"
	"sync"
)

type PlaylistInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	MediaID string `json:"mediaId"`
}

var (
	playlists     []PlaylistInfo
	playlistsOnce sync.Once
	playlistsMu   sync.Mutex
)

func loadPlaylists() error {
	playlistsOnce.Do(func() {
		playlists = make([]PlaylistInfo, 0)
		f, err := os.Open("playlists.dat")
		if os.IsNotExist(err) {
			return
		}
		if err != nil {
			return
		}
		defer f.Close()
		json.NewDecoder(f).Decode(&playlists)
	})
	return nil
}

func getPlaylists() ([]PlaylistInfo, error) {
	playlistsMu.Lock()
	defer playlistsMu.Unlock()

	if err := loadPlaylists(); err != nil {
		return nil, err
	}
	return playlists, nil
}

func createPlaylist(name string, mediaId string) (string, error) {
	playlistsMu.Lock()
	defer playlistsMu.Unlock()

	if err := loadPlaylists(); err != nil {
		return "", err
	}

	newPlaylist := PlaylistInfo{
		ID:      "bili-" + mediaId,
		Name:    name,
		MediaID: mediaId,
	}

	playlists = append(playlists, newPlaylist)

	f, err := os.Create("playlists.dat")
	if err != nil {
		return "", err
	}
	defer f.Close()

	return newPlaylist.ID, json.NewEncoder(f).Encode(playlists)
}
