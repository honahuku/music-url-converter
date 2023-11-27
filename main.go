package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func main() {
	// YouTube MusicのURL
	youtubeMusicURL := "https://music.youtube.com/watch?v=EVBsypHzF3U&si=HuSCNAk7ZFbkOAB5"

	// 正規表現でビデオIDを抽出
	re := regexp.MustCompile(`(?<=watch\?v=)[a-zA-Z0-9]+(?=&|$)`)
	videoID := re.FindString(youtubeMusicURL)
	if videoID == "" {
		fmt.Println("No video ID found in URL")
		return
	}

	// 環境変数からYouTube APIキーを取得
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		fmt.Println("YOUTUBE_API_KEY is not set")
		return
	}

	// YouTube APIのURLを構築
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/videos?part=snippet,contentDetails,statistics&id=%s&key=%s", videoID, apiKey)

	// HTTPリクエストを作成
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// レスポンスのボディを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// レスポンスを出力
	fmt.Println(string(body))
}
