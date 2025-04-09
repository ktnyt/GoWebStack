package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/websocket"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func reloader(wg *sync.WaitGroup, sigCh <-chan os.Signal) {
	defer wg.Done()

	e := echo.New()

	targets := make(map[string]chan bool)

	e.GET("/subscribe", func(c echo.Context) error {
		fmt.Printf("%s[RELOADER] ブラウザに接続します...%s\n", colorBlue, colorReset)
		id := uuid.New().String()
		targets[id] = make(chan bool)

		websocket.Handler(func(ws *websocket.Conn) {
			defer func() {
				close(targets[id])
				delete(targets, id)
				ws.Close()
			}()

			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-targets[id]:
					fmt.Printf("%s[RELOADER] ブラウザにリロードを通知します...%s\n", colorBlue, colorReset)
					ws.Write([]byte("reload"))
				case <-ticker.C:
					ws.Write([]byte("ping"))
				}
			}
		}).ServeHTTP(c.Response().Writer, c.Request())

		return nil
	})

	e.GET("/reload", func(c echo.Context) error {
		fmt.Printf("%s[RELOADER] リロードを開始します...%s\n", colorBlue, colorReset)

		client := http.Client{Timeout: 1 * time.Second}
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				resp, err := client.Get("http://localhost:6640/status")
				if err != nil {
					fmt.Printf("%s[RELOADER] サーバーに接続できませんでした...%s\n", colorBlue, colorReset)
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					fmt.Printf("%s[RELOADER] サーバーの起動が確認できました...%s\n", colorBlue, colorReset)
					for _, target := range targets {
						target <- true
					}
					return nil
				}
			}
		}
	})

	// サーバーを非同期で起動
	go func() {
		if err := e.Start(":6641"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("サーバー起動エラー: %v", err)
		}
	}()

	fmt.Printf("%s[RELOADER] サーバーを起動しました%s\n", colorBlue, colorReset)

	// シグナルを待機
	<-sigCh
	fmt.Printf("%s[RELOADER] シャットダウンを開始します...%s\n", colorBlue, colorReset)

	// graceful shutdownのためのコンテキスト
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// サーバーをシャットダウン
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Printf("%s[RELOADER] サーバーを正常に終了しました%s\n", colorBlue, colorReset)
}
