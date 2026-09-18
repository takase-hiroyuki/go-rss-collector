package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
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

	// 追記モード（無ければ新規作成）で articles.md を開く
	file, err := os.OpenFile("articles.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("ファイルオープンエラー: %v\n", err)
		return
	}
	defer file.Close()

	// 実行日時のヘッダーを書き込む
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(file, "\n## 取得日時: %s\n\n", now)

	// 最新5件をMarkdownの箇条書きで書き込む
	count := 0
	for _, item := range rdf.Items {
		if count >= 5 {
			break
		}
		fmt.Fprintf(file, "* [%s](%s)\n", item.Title, item.Link)
		count++
	}

	fmt.Println("articles.md への書き込みが完了しました。")
}
