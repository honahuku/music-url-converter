package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

// getToken retrieves an access token from the Spotify API using the provided
// client ID and client secret. It returns the token and any error encountered.
func getToken(clientID, clientSecret string) (string, error) {
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
func searchSpotify(query, token string) (string, error) {
	encodedQuery := url.QueryEscape(query)
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", encodedQuery)

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

func main() {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET must be set")
	}

	token, err := getToken(clientID, clientSecret)
	if err != nil {
		log.Fatalf("Error getting token: %v", err)
	}

	query := "藤井風 きらり"
	result, err := searchSpotify(query, token)
	if err != nil {
		log.Fatalf("Error searching Spotify: %v", err)
	}

	// 結果からトラックのURLを抽出します。
	trackURL, err := extractTrackURL(result)
	if err != nil {
		log.Fatalf("Error extracting track URL: %v", err)
	}

	fmt.Println("Track URL:", trackURL)
}