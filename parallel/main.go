package main

import (
	"fmt"
	"sync"
	"time"
)

// 模擬的に時間がかかる処理を行う関数
func fetchTask(id int, name string, wg *sync.WaitGroup) {
	defer wg.Done() // 関数終了時に WaitGroup のカウントを1減らす

	fmt.Printf("[開始] タスク %d: %s\n", id, name)
	// 通信待ちを想定した1秒間の待機
	time.Sleep(1 * time.Second)
	fmt.Printf("  -> [完了] タスク %d: %s\n", id, name)
}

func main() {
	var wg sync.WaitGroup

	tasks := []string{
		"はてなブックマークの取得",
		"Qiitaトレンドの取得",
		"Zenn新着記事の取得",
	}

	start := time.Now()
	fmt.Println("=== 並行処理を開始します ===")

	for i, task := range tasks {
		wg.Add(1) // 待機するタスク数を1増やす
		// 関数の呼び出し前に「go」を付けるだけで並行処理（Goroutine）が起動
		go fetchTask(i+1, task, &wg)
	}

	// すべてのタスクの wg.Done() が呼ばれるまで待機
	wg.Wait()

	fmt.Printf("=== すべて完了（総所要時間: %v）===\n", time.Since(start).Round(time.Millisecond))
}
