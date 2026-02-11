# go-test-splitting-sample

`Goテストファイル分割の実践` 記事向けの最小サンプルです。

## 構成

- `user/user_service.go`
- `user/user_service_unit_test.go`
- `user/user_service_integration_test.go`
- `user/user_service_scenario_test.go`

## 実行例

```bash
go test ./...
```

`integration` タグ付きテストを含める場合:

```bash
go test ./... -tags=integration
```
