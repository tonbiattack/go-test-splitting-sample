package user

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrInvalidUserID = errors.New("invalid user id")
	ErrInvalidSecret = errors.New("invalid secret")
	ErrUserNotFound  = errors.New("user not found")
	ErrUserLocked    = errors.New("user locked")
	ErrAuthFailed    = errors.New("authentication failed")
)

const maxFailedLoginCount = 3

type User struct {
	ID               string
	FailedLoginCount int
	Locked           bool
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return UserService{repo: repo}
}

type UserRepository interface {
	Exists(ctx context.Context, id string) (bool, error)
	FindByID(ctx context.Context, id string) (User, error)
	IncrementFailedLogin(ctx context.Context, id string) (User, error)
	ResetFailedLogin(ctx context.Context, id string) error
	AppendAuditLog(ctx context.Context, id string, result string) error
}

func (s UserService) CanLogin(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, ErrInvalidUserID
	}

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return false, fmt.Errorf("check exists: %w", err)
	}
	if !exists {
		return false, ErrUserNotFound
	}

	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("find user: %w", err)
	}
	if u.Locked {
		return false, ErrUserLocked
	}

	return true, nil
}

type PasswordVerifier interface {
	Verify(ctx context.Context, id string, secret string) (bool, error)
}

func (s UserService) Login(
	ctx context.Context,
	id string,
	secret string,
	verifier PasswordVerifier,
) error {
	if secret == "" {
		return ErrInvalidSecret
	}

	canLogin, err := s.CanLogin(ctx, id)
	if err != nil {
		return err
	}
	if !canLogin {
		return ErrUserLocked
	}

	matched, err := verifier.Verify(ctx, id, secret)
	if err != nil {
		return fmt.Errorf("verify secret: %w", err)
	}
	if !matched {
		u, incErr := s.repo.IncrementFailedLogin(ctx, id)
		if incErr != nil {
			return fmt.Errorf("increment failed login: %w", incErr)
		}
		if auditErr := s.repo.AppendAuditLog(ctx, id, "failed"); auditErr != nil {
			return fmt.Errorf("append audit log: %w", auditErr)
		}
		if u.FailedLoginCount >= maxFailedLoginCount {
			return ErrUserLocked
		}
		return ErrAuthFailed
	}

	if err := s.repo.ResetFailedLogin(ctx, id); err != nil {
		return fmt.Errorf("reset failed login: %w", err)
	}
	if err := s.repo.AppendAuditLog(ctx, id, "success"); err != nil {
		return fmt.Errorf("append audit log: %w", err)
	}

	return nil
}
