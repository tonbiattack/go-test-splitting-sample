//go:build integration
// +build integration

package user

import "testing"

func TestUserService_Integration_CanLoginWithDB(t *testing.T) {
	// テストDBを起動し、Repository実装を挿して確認する想定。
	t.Skip("example integration test")
}
