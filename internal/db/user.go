package db

import "go-server-template/internal/model"

func GetUserByID(id uint) (user *model.User, err error) {
	user = &model.User{ID: id}
	if err = db.Debug().First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
