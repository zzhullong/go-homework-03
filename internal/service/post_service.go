package service

import (
	"homework03/internal/dao"
	"homework03/internal/model"
	"homework03/pkg/db"
)

type PostService struct {
	postDao 	*dao.PostDao
	userDao    *dao.UserDao
	
}

func NewPostService() *PostService {
	return &PostService{
		postDao: dao.NewPostDao(),
		userDao: dao.NewUserDao(),
	}
}

func (s *PostService) GetMostCommentPost() (*model.Post, error) {
	return s.postDao.GetMostCommentPost()
}

func (s *PostService) CreatePost(title, content string, userId uint) (*model.Post, error) {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback().Error
		}
	}()

	post := &model.Post{
		Title:   title,
		Content: content,
		UserID:  userId,
	}
	if err := s.postDao.CreatePost(tx, post); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetPostByID(pid uint) (*model.Post, error) {
	return s.postDao.GetPostByID(pid)
}

func (s *PostService) DeletePost(postId uint) error {
	return s.postDao.DeletePost(postId)
}