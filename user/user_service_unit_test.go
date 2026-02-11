package user

import (
	"context"
	"errors"
	"testing"
)

type stubRepo struct {
	exists       bool
	existsErr    error
	user         User
	findErr      error
	incremented  User
	incrementErr error
	resetErr     error
	auditErr     error
	auditLogs    []string
}

func (r stubRepo) Exists(context.Context, string) (bool, error) {
	return r.exists, r.existsErr
}

func (r stubRepo) FindByID(context.Context, string) (User, error) {
	return r.user, r.findErr
}

func (r stubRepo) IncrementFailedLogin(context.Context, string) (User, error) {
	return r.incremented, r.incrementErr
}

func (r stubRepo) ResetFailedLogin(context.Context, string) error {
	return r.resetErr
}

func (r *stubRepo) AppendAuditLog(context.Context, string, string) error {
	if r.auditErr != nil {
		return r.auditErr
	}
	r.auditLogs = append(r.auditLogs, "logged")
	return nil
}

type stubVerifier struct {
	matched bool
	err     error
}

func (v stubVerifier) Verify(context.Context, string, string) (bool, error) {
	return v.matched, v.err
}

func TestUserService_Unit_CanLogin(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		repo    stubRepo
		wantOK  bool
		wantErr error
	}{
		{
			name: "success",
			id:   "u1",
			repo: stubRepo{
				exists: true,
				user:   User{ID: "u1", Locked: false},
			},
			wantOK: true,
		},
		{
			name:    "invalid id",
			id:      "",
			repo:    stubRepo{},
			wantOK:  false,
			wantErr: ErrInvalidUserID,
		},
		{
			name: "not found",
			id:   "unknown",
			repo: stubRepo{
				exists: false,
			},
			wantOK:  false,
			wantErr: ErrUserNotFound,
		},
		{
			name: "locked user",
			id:   "u2",
			repo: stubRepo{
				exists: true,
				user:   User{ID: "u2", Locked: true},
			},
			wantOK:  false,
			wantErr: ErrUserLocked,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.repo
			svc := NewUserService(&repo)
			ok, err := svc.CanLogin(context.Background(), tc.id)

			if ok != tc.wantOK {
				t.Fatalf("ok mismatch: got=%v want=%v", ok, tc.wantOK)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err mismatch: got=%v want=%v", err, tc.wantErr)
			}
		})
	}
}

func TestUserService_Unit_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &stubRepo{
			exists: true,
			user:   User{ID: "u1"},
		}
		svc := NewUserService(repo)

		err := svc.Login(context.Background(), "u1", "secret", stubVerifier{matched: true})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(repo.auditLogs) != 1 {
			t.Fatalf("expected audit log count=1, got %d", len(repo.auditLogs))
		}
	})

	t.Run("wrong secret increments failed count", func(t *testing.T) {
		repo := &stubRepo{
			exists:      true,
			user:        User{ID: "u1"},
			incremented: User{ID: "u1", FailedLoginCount: 1},
		}
		svc := NewUserService(repo)

		err := svc.Login(context.Background(), "u1", "wrong", stubVerifier{matched: false})
		if !errors.Is(err, ErrAuthFailed) {
			t.Fatalf("expected ErrAuthFailed, got %v", err)
		}
		if len(repo.auditLogs) != 1 {
			t.Fatalf("expected audit log count=1, got %d", len(repo.auditLogs))
		}
	})

	t.Run("third failure returns lock error", func(t *testing.T) {
		repo := &stubRepo{
			exists:      true,
			user:        User{ID: "u1"},
			incremented: User{ID: "u1", FailedLoginCount: 3, Locked: true},
		}
		svc := NewUserService(repo)

		err := svc.Login(context.Background(), "u1", "wrong", stubVerifier{matched: false})
		if !errors.Is(err, ErrUserLocked) {
			t.Fatalf("expected ErrUserLocked, got %v", err)
		}
	})
}
