package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

// getSpotifyToken retrieves an access token from the Spotify API using the provided
// client ID and client secret. It returns the token and any error encountered.
func getSpotifyToken(clientID, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

// searchSpotify searches the Spotify database for tracks matching the given query.
// It returns the raw JSON response and any error encountered.
func searchSpotify(query, token string, searchType string) (string, error) {
	encodedQuery := url.QueryEscape(query)
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=%s", encodedQuery, searchType)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// extractTrackURL parses the provided JSON data to extract the Spotify URL of the first track.
// It returns the URL as a string and any error encountered. If no tracks are found, it returns an error.
func extractTrackURL(jsonData string) (string, error) {
	var result struct {
		Tracks struct {
			Items []struct {
				ExternalURLs struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
			} `json:"items"`
		} `json:"tracks"`
	}

	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return "", err
	}

	if len(result.Tracks.Items) > 0 {
		return result.Tracks.Items[0].ExternalURLs.Spotify, nil
	}

	return "", fmt.Errorf("no tracks found")
}

// performSpotifySearch performs a Spotify search for the given query and search type.
func performSpotifySearch(query, searchType string) (string, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET must be set")
	}

	token, err := getSpotifyToken(clientID, clientSecret)
	if err != nil {
		return "", fmt.Errorf("Error getting Spotify token: %v", err)
	}

	result, err := searchSpotify(query, token, searchType)
	if err != nil {
		return "", fmt.Errorf("Error searching Spotify: %v", err)
	}

	return extractTrackURL(result)
}

// TrackInfo is a struct for storing track information.
type TrackInfo struct {
	TrackName  string `json:"trackName"`
	ArtistName string `json:"artistName"`
}

// getYoutubeMusicInfo gets track information from YouTube Music.
func getYoutubeMusicInfo(url string) (TrackInfo, error) {
	// TODO: #3 Youtube Music以外にも対応する
	// Node.jsのスクリプトを実行
	dir := "./downloadPage"

	buildCmd := exec.Command("npm", "run", "build")
	buildCmd.Dir = dir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	buildCmd.Stdout = &stdout
	buildCmd.Stderr = &stderr
	err := buildCmd.Run()
	if err != nil {
		// コマンド実行に失敗した場合の処理
		fmt.Printf("Error running npm run build: %v\n", err)
		if stderr.String() != "" {
			fmt.Printf("Stderr: %s\n", stderr.String())
		}
	} else {
		// コマンド実行に成功した場合の処理
		if stdout.String() != "" {
			fmt.Printf("npm run build success. Stdout: %s", stdout.String())
		} else {
			fmt.Printf("npm run build success but no output.")
		}
	}

	cmd := exec.Command("npm", "start", "--", url)
	cmd.Dir = dir
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return TrackInfo{}, err
	}

	if err := cmd.Start(); err != nil {
		return TrackInfo{}, err
	}

	scanner := bufio.NewScanner(stdoutPipe)
	var info TrackInfo
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line) // リアルタイム出力
		// JSON形式の行を検出した場合にデコード
		if json.Valid([]byte(line)) {
			if err := json.Unmarshal([]byte(line), &info); err != nil {
				log.Println("Error decoding JSON:", err)
				continue
			}
			break // JSONをデコードしたらループを抜ける
		}
	}

	if err := cmd.Wait(); err != nil {
		return TrackInfo{}, err
	}

	return info, nil
}

func main() {
	// YouTube MusicのURL
	youtubeMusicURL := "https://music.youtube.com/watch?v=VNQB91Ra_rs&si=7zmrGPXhE0WdrPCM"

	// YouTube Musicからトラック情報を取得
	info, err := getYoutubeMusicInfo(youtubeMusicURL)
	if err != nil {
		log.Fatalf("Error getting YouTube Music info: %v", err)
	}

	query := fmt.Sprintf("%s %s", info.ArtistName, info.TrackName)
	searchType := "track" // TODO: #1 track以外の検索もできるようにする

	trackURL, err := performSpotifySearch(query, searchType)
	if err != nil {
		log.Fatalf("Error searching Spotify: %v", err)
	}
	fmt.Println("Track URL:", trackURL)
}
