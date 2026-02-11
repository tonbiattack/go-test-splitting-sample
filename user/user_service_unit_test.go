package user

import (
	"context"
	"testing"
)

type stubRepo struct {
	exists bool
	err    error
}

func (r stubRepo) Exists(context.Context, string) (bool, error) {
	return r.exists, r.err
}

func TestUserService_Unit_CanLogin(t *testing.T) {
	svc := NewUserService(stubRepo{exists: true})
	ok, err := svc.CanLogin(context.Background(), "u1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true")
	}
}
