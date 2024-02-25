package service

import "gorm.io/gorm"

var svc *service

func Init(db *gorm.DB) {
	svc = &service{db: db}
}

func Get() Service {
	return svc
}

var _ Service = (*service)(nil)

type Service interface {
	User() UserService

	i()
}
type service struct {
	db *gorm.DB
}

func (s *service) User() UserService {
	return newUser(s)
}

func (s *service) i() {}
