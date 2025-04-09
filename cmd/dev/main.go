package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// シグナルを待ち受けるチャネル
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// チャネルを作成してコマンドを追跡
	cmds := make(chan *exec.Cmd, 3)
	var wg sync.WaitGroup
	wg.Add(3)

	// 各プロセスを開始
	go reloader(&wg, sigCh)
	go tailwind(&wg, cmds)
	go templ(&wg, cmds)
	go server(&wg, cmds)

	// 子プロセスを追跡
	processes := []*exec.Cmd{}
	go func() {
		for cmd := range cmds {
			processes = append(processes, cmd)
		}
	}()

	fmt.Println("開発サーバーを起動しました。Ctrl+C で終了します。")

	// シグナルを受け取ったら子プロセスを停止
	<-sigCh
	fmt.Println("\n終了シグナルを受信しました。すべてのプロセスを停止しています...")

	for _, proc := range processes {
		if proc != nil && proc.Process != nil {
			proc.Process.Kill()
		}
	}

	// 子プロセスの終了を待つ
	close(cmds)
	wg.Wait()
	fmt.Println("すべてのプロセスが終了しました。")
}
