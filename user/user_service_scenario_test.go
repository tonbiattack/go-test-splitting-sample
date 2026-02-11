package user

import (
	"context"
	"errors"
	"testing"
)

func TestUserService_Scenario_LoginFlow(t *testing.T) {
	repo := NewInMemoryUserRepository([]User{
		{ID: "u1"},
	})
	verifier := NewStaticPasswordVerifier(map[string]string{
		"u1": "correct",
	})
	svc := NewUserService(repo)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		err := svc.Login(ctx, "u1", "wrong", verifier)
		if !errors.Is(err, ErrAuthFailed) {
			t.Fatalf("attempt=%d expected ErrAuthFailed, got %v", i+1, err)
		}
	}

	err := svc.Login(ctx, "u1", "wrong", verifier)
	if !errors.Is(err, ErrUserLocked) {
		t.Fatalf("third failure should lock user, got %v", err)
	}

	err = svc.Login(ctx, "u1", "correct", verifier)
	if !errors.Is(err, ErrUserLocked) {
		t.Fatalf("locked user should not login, got %v", err)
	}

	logs := repo.AuditLog("u1")
	if len(logs) != 3 {
		t.Fatalf("expected 3 failed audit logs, got %d", len(logs))
	}
}
