package main

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

const starredFile = "starred.dat"

var mu sync.Mutex

// getStarredSongs_nl reads the starred songs file without locking.
func getStarredSongs_nl() ([]string, error) {
	f, err := os.Open(starredFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var songs []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		songs = append(songs, scanner.Text())
	}
	return songs, scanner.Err()
}

// starSong adds a song ID to the starred file.
func starSong(id string) error {
	mu.Lock()
	defer mu.Unlock()

	songs, err := getStarredSongs_nl()
	if err != nil {
		return err
	}

	for _, songID := range songs {
		if songID == id {
			return nil // Already starred
		}
	}

	f, err := os.OpenFile(starredFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(id + "\n")
	return err
}

// unstarSong removes a song ID from the starred file.
func unstarSong(id string) error {
	mu.Lock()
	defer mu.Unlock()

	songs, err := getStarredSongs_nl()
	if err != nil {
		return err
	}

	var newSongs []string
	for _, songID := range songs {
		if songID != id {
			newSongs = append(newSongs, songID)
		}
	}

	return os.WriteFile(starredFile, []byte(strings.Join(newSongs, "\n")+"\n"), 0644)
}

// getStarredSongs returns a list of starred song IDs.
func getStarredSongs() ([]string, error) {
	mu.Lock()
	defer mu.Unlock()
	return getStarredSongs_nl()
}