# go-test-splitting-sample

`Goテストファイル分割の実践` 記事向けのサンプルです。
1つの `user_service.go` に対して、`unit / integration / scenario` の3種類でテストを分けています。

## 仕様

- `UserService.CanLogin`
  - IDの妥当性
  - ユーザー存在確認
  - ロック状態判定
- `UserService.Login`
  - パスワード検証
  - 認証失敗回数の加算
  - 3回失敗でロック
  - 監査ログ記録

## ファイル構成

- `user/user_service.go`
  - ユースケース本体
- `user/inmemory_repository.go`
  - 疑似DBとして使うインメモリRepository
- `user/static_password_verifier.go`
  - 固定値で検証する認証器
- `user/user_service_unit_test.go`
  - スタブでサービスロジックを検証
- `user/user_service_integration_test.go`
  - `integration` タグ付き。Repository実装込みで状態変化を検証
- `user/user_service_scenario_test.go`
  - ログイン失敗からロックまでの業務フローを検証

## 実行コマンド

通常実行（unit + scenario）:

```bash
go test ./...
```

integrationを含めて実行:

```bash
go test ./... -tags=integration
```

## なぜ `-tags=integration` が必要か

`user/user_service_integration_test.go` には次の build tag を付けています。

```go
//go:build integration
// +build integration
```

この指定があるファイルは、`go test` 実行時に `-tags=integration` を付けた場合だけ
テスト対象に含まれます。タグなしだとファイル自体が除外されるため、
`-run TestUserService_Integration...` を指定しても `no tests to run` になります。
