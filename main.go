package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// RDF / RSS 1.0 に対応した構造体
type RDF struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

func main() {
	url := "https://b.hatena.ne.jp/hotentry/it.rss"

	// 1. RSS フィードの取得
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("取得エラー: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var rdf RDF
	decoder := xml.NewDecoder(resp.Body)
	if err := decoder.Decode(&rdf); err != nil {
		fmt.Printf("パースエラー: %v\n", err)
		return
	}

	// 2. 既存の articles.md の中身を読み込む（既読チェック用）
	existingContent := ""
	data, err := os.ReadFile("articles.md")
	if err == nil {
		existingContent = string(data)
	}

	// 3. 追記モードで articles.md を開く
	file, err := os.OpenFile("articles.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("ファイルオープンエラー: %v\n", err)
		return
	}
	defer file.Close()

	// 4. 日本時間（JST）を取得して日時文字列を整形
	jst := time.FixedZone("JST", 9*60*60)
	now := time.Now().In(jst)
	dateStr := fmt.Sprintf("%d月%d日 %d時%d分", now.Month(), now.Day(), now.Hour(), now.Minute())

	// 5. 新着記事（既存ファイルに URL が含まれていないもの）を抽出
	var newItems []Item
	for _, item := range rdf.Items {
		if !strings.Contains(existingContent, item.Link) {
			newItems = append(newItems, item)
		}
		if len(newItems) >= 5 { // 最大5件まで
			break
		}
	}

	// 6. 判定とファイルへの書き込み
	if len(newItems) == 0 {
		// 新規記事がない場合
		fmt.Fprintf(file, "\n%s 変更ありません\n", dateStr)
		fmt.Printf("%s 変更ありません（追記完了）\n", dateStr)
	} else {
		// 新規記事がある場合
		fmt.Fprintf(file, "\n## 取得日時: %s\n\n", dateStr)
		for _, item := range newItems {
			fmt.Fprintf(file, "* [%s](%s)\n", item.Title, item.Link)
		}
		fmt.Printf("%d 件の新規記事を追記しました。\n", len(newItems))
	}
}
