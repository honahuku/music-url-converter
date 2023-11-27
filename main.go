package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// YouTube APIレスポンスのための構造体
type YouTubeResponse struct {
	Items []struct {
		Snippet struct {
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}

func main() {
	// TODO: ytmusicのurlからidを取得するまでの処理を関数に切り出す

	// YouTube MusicのURL
	youtubeMusicURL := "https://music.youtube.com/watch?v=EVBsypHzF3U&si=HuSCNAk7ZFbkOAB5"
	// "watch?v=" の後ろを取得
	idStartIndex := strings.Index(youtubeMusicURL, "watch?v=") + len("watch?v=")
	if idStartIndex == -1 {
		fmt.Println("Invalid URL format")
		return
	}

	// "&si=" がある場合、それ以降を削除
	videoID := youtubeMusicURL[idStartIndex:]
	if siIndex := strings.Index(videoID, "&si="); siIndex != -1 {
		videoID = videoID[:siIndex]
	}

	// TODO: videoIDを引数にしてタイトルを返却する、YouTube APIを呼び出す処理を関数に切り出す

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

	// レスポンスをJSONとして解析
	var ytResp YouTubeResponse
	if err := json.Unmarshal(body, &ytResp); err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return
	}

	// items[0].snippet.titleを出力
	if len(ytResp.Items) > 0 {
		fmt.Println("Title:", ytResp.Items[0].Snippet.Title)
	} else {
		fmt.Println("No items found in response")
	}
}
