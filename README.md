# servicex

go でバックエンドサーバーを建てるときに必要な共通ライブラリ

## servicex の担っていること

- 環境変数(`.env`等)関連の処理
- [Sentry](https://sentry.io/)にエラーを送る
- log を Google Cloud Logging に送る

## 使い方

### 新しく作るサービスに書くべきコード

```go:main.go
package main

import (
  "log"
  "github.com/rimoapp/servicex"
  "github.com/rimoapp/yourapp/server" // code is in your app
)

func main() {
  servicex.Init()
	defer servicex.Close()

  // run server
  if err := server.Run(); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}
```

```go:server/router.go
package server

import (
  "github.com/rimoapp/servicex"
)

// Run runs server on :8080 unless a PORT environment variable was defined.
func Run() error {
  router := newRouter()
  return router.Run()
}

func newRouter() *gin.Engine {
  router := servicex.DefaultGinEngine()

  router.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  }
  // ...
}
```

### 設定するべき環境変数

- `SENTRY_DNS`
  - Sentry にて、対応する新規プロジェクトをリポジトリと同名で作成して取得
- `GOOGLE_CLOUD_PROJECT`
  - Cloud Logging に送るための Google Cloud のプロジェクトを設定をする(e.g. `rimo-dev-0`, `rimo-prod`)
- `APP_ENV`
  - 利用環境 (e.g. `development`, `staging`, `production`)
- `GIT_COMMIT_SHA` (optional)
  - Sentry にどのリリースなのかレポートされる。Github Actions で `github.sha` を引き渡すことにより設定
- `PORT` (optional)
  - 設定するとポートを変えられる。デフォルトは`:8080`
