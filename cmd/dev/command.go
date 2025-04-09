package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

// 色のコード
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

// プロセスの出力を受け取って色付きで表示する
func pipeOutput(prefix string, color string, r io.Reader, isError bool) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		if isError {
			fmt.Fprintf(os.Stderr, "%s[%s] %s%s\n", color, prefix, text, colorReset)
		} else {
			fmt.Printf("%s[%s] %s%s\n", color, prefix, text, colorReset)
		}
	}
}

func runCommand(name string, color string, args []string, wg *sync.WaitGroup, cmds chan<- *exec.Cmd, useTTY bool) {
	defer wg.Done()

	cmd := exec.Command(args[0], args[1:]...)

	// TTYモードの設定
	if useTTY {
		// 現在のプロセスの標準入力を子プロセスにも設定
		cmd.Stdin = os.Stdin
		// 環境変数でフォースカラーモードを有効化
		cmd.Env = append(os.Environ(),
			"FORCE_COLOR=true",
			"COLORTERM=truecolor",
			"TERM=xterm-256color")
	}

	// 標準出力パイプの設定
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s[%s] Failed to create stdout pipe: %v%s\n", color, name, err, colorReset)
		return
	}

	// 標準エラーパイプの設定
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s[%s] Failed to create stderr pipe: %v%s\n", color, name, err, colorReset)
		return
	}

	// 出力をリアルタイムで処理
	go pipeOutput(name, color, stdout, false)
	go pipeOutput(name, color, stderr, true)

	// コマンドチャンネルに追加
	cmds <- cmd

	// コマンド開始
	fmt.Printf("%s[%s] Starting process%s\n", color, name, colorReset)
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "%s[%s] Failed to start: %v%s\n", color, name, err, colorReset)
		return
	}

	// コマンド終了を待機
	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "%s[%s] Process ended with error: %v%s\n", color, name, err, colorReset)
	} else {
		fmt.Printf("%s[%s] Process completed successfully%s\n", color, name, colorReset)
	}
}
