package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

// getSpotifyClient はSpotify APIのクライアントを生成して返します。
func getSpotifyClient(ctx context.Context) *spotify.Client {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatalf("SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET must be set")
	}

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	return client
}

// getTrackURL は与えられたアーティスト名と曲名に基づいてSpotifyでトラックを検索し、
// 見つかった最初のトラックのURLを返します。
func getTrackURL(client *spotify.Client, artistName string, trackName string) (string, error) {
	ctx := context.Background()

	// URLエンコードされた検索クエリを作成します。
	query := url.QueryEscape(fmt.Sprintf("artist:%s track:%s", artistName, trackName))

	// エンコードされたクエリを使って検索を行います。
	searchResults, err := client.Search(ctx, query, spotify.SearchTypeTrack)
	if err != nil {
		fmt.Printf("Search error: %v\n", err)
		return "", err
	}

	// 検索結果から最初のトラックを取得します。
	if searchResults.Tracks != nil && len(searchResults.Tracks.Tracks) > 0 {
		track := searchResults.Tracks.Tracks[0]
		return track.ExternalURLs["spotify"], nil
	}

	return "", fmt.Errorf("no track found")
}

func main() {
	ctx := context.Background()
	client := getSpotifyClient(ctx)

	// ここでアーティスト名と曲名を指定します。
	artistName := "kaze fujii"
	trackName := "死ぬのがいいわ"

	// トラックのURLを取得します。
	trackURL, err := getTrackURL(client, artistName, trackName)
	if err != nil {
		log.Fatalf("Error getting track URL: %v", err)
	}

	fmt.Println("Track URL:", trackURL)
}
