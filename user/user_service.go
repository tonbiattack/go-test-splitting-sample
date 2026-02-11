package user

import "context"

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return UserService{repo: repo}
}

type UserRepository interface {
	Exists(ctx context.Context, id string) (bool, error)
}

func (s UserService) CanLogin(ctx context.Context, id string) (bool, error) {
	return s.repo.Exists(ctx, id)
}
