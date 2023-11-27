package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	// YouTube MusicのURL
	youtubeMusicURL := "https://music.youtube.com/watch?v=SepkbjeTe7I&si=mG8D135LUv82zR7p"

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
