package service

import (
	"context"
	"go-server-template/internal/db"
	"go-server-template/internal/model"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uint) (user *model.User, err error)

	i()
}

type userService struct {
	db *gorm.DB
}

func newUser(s *service) UserService {
	return &userService{
		db: s.db,
	}
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (user *model.User, err error) {
	return db.GetUserByID(ctx, id)
}

func (s *userService) i() {}
