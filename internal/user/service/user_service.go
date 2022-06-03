package service

import (
	"context"
	"database/sql"
	q "github.com/core-go/sql"

	. "go-service/internal/user/domain"
	. "go-service/internal/user/port"
)

type UserService interface {
	Load(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) (int64, error)
	Update(ctx context.Context, user *User) (int64, error)
	Patch(ctx context.Context, user map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}

func NewUserService(db *sql.DB, repository UserRepository) UserService {
	return &userService{db: db, repository: repository}
}

type userService struct {
	db         *sql.DB
	repository UserRepository
}

func (s *userService) Load(ctx context.Context, id string) (*User, error) {
	return s.repository.Load(ctx, id)
}
func (s *userService) Create(ctx context.Context, user *User) (int64, error) {
	ctx, tx, err := q.Begin(ctx, s.db)
	if err != nil {
		return  -1, err
	}
	res, err := s.repository.Create(ctx, user)
	return q.End(tx, res, err)
}
func (s *userService) Update(ctx context.Context, user *User) (int64, error) {
	ctx, tx, err := q.Begin(ctx, s.db)
	if err != nil {
		return  -1, err
	}
	res, err := s.repository.Update(ctx, user)
	err = q.Commit(tx, err)
	return res, err
}
func (s *userService) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	ctx, tx, err := q.Begin(ctx, s.db)
	if err != nil {
		return  -1, err
	}
	res, err := s.repository.Patch(ctx, user)
	err = q.Commit(tx, err)
	return res, err
}
func (s *userService) Delete(ctx context.Context, id string) (int64, error) {
	ctx, tx, err := q.Begin(ctx, s.db)
	if err != nil {
		return  -1, err
	}
	res, err := s.repository.Delete(ctx, id)
	err = q.Commit(tx, err)
	return res, err
}
