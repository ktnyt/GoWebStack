package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// routes.goで定義したルーティングをセットアップ
	SetupRoutes(e)

	// シグナルを待ち受けるチャネル
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// サーバーを非同期で起動
	go func() {
		if err := e.Start(":6640"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("サーバー起動エラー: %v", err)
		}
	}()

	// シグナルを待機
	<-quit
	e.Logger.Info("シャットダウンを開始します...")

	// graceful shutdownのためのコンテキスト
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// サーバーをシャットダウン
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Info("サーバーを正常に終了しました")
}
