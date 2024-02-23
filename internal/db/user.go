package db

import (
	"context"
	"go-server-template/internal/model"
)

func GetUserByID(ctx context.Context, id uint) (user *model.User, err error) {
	user = &model.User{ID: id}
	if err = db.WithContext(ctx).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
