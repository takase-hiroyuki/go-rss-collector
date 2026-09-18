package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// RSS 2.0 の構造定義
type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

func main() {
	// 例: はてなブックマーク テクノロジーのRSS 2.0
	url := "https://b.hatena.ne.jp/hotentry/it.rss"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("取得エラー: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var rss RSS
	decoder := xml.NewDecoder(resp.Body)
	if err := decoder.Decode(&rss); err != nil {
		fmt.Printf("パースエラー: %v\n", err)
		return
	}

	fmt.Printf("=== 取得成功: %s ===\n", rss.Channel.Title)
	for i, item := range rss.Channel.Items {
		if i >= 5 { // 最新5件を表示
			break
		}
		fmt.Printf("[%d] %s\n    %s\n", i+1, item.Title, item.Link)
	}
}
