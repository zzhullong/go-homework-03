package service

import (
	"homework03/internal/dao"
	"homework03/internal/model"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService() *UserService {
	return &UserService{
		userDao: dao.NewUserDao(),
	}
}

func (s *UserService) GetUserWithPostsAndComments(uid int) (*model.User, error) {
	return s.userDao.GetUserWithPostsAndComments(uid)
}

func (s *UserService) CreateUser(username, password string) (*model.User, error) {
	user := &model.User{
		Name: username,
		Pwd: password,
	}
	err := s.userDao.CreateUser(user)
	return user, err
}

func (s *UserService) DeleteUser(userId int) error {
	return s.userDao.DeleteUser(userId)
}
