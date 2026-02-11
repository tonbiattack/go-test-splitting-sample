//go:build integration
// +build integration

package user

import (
	"context"
	"errors"
	"testing"
)

func TestUserService_Integration_CanLoginWithDB(t *testing.T) {
	repo := NewInMemoryUserRepository([]User{
		{ID: "u1"},
	})
	verifier := NewStaticPasswordVerifier(map[string]string{
		"u1": "correct",
	})
	svc := NewUserService(repo)

	err := svc.Login(context.Background(), "u1", "wrong", verifier)
	if !errors.Is(err, ErrAuthFailed) {
		t.Fatalf("expected ErrAuthFailed, got %v", err)
	}

	u, err := repo.FindByID(context.Background(), "u1")
	if err != nil {
		t.Fatalf("expected find success, got %v", err)
	}
	if u.FailedLoginCount != 1 {
		t.Fatalf("expected failed_login_count=1, got %d", u.FailedLoginCount)
	}

	logs := repo.AuditLog("u1")
	if len(logs) != 1 || logs[0] != "failed" {
		t.Fatalf("unexpected audit logs: %#v", logs)
	}
}
