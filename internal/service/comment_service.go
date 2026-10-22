package service

import (
	"homework03/internal/dao"
	"homework03/internal/model"
	"homework03/pkg/db"
)

type CommentService struct {
	commentDao *dao.CommentDao
	postDao *dao.PostDao
}

func NewCommentService() *CommentService {
	return &CommentService{
		commentDao: dao.NewCommentDao(),
		postDao: dao.NewPostDao(),
	}
}

func (s *CommentService) CreateComment(content string, postId uint) (*model.Comment, error) {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback().Error
		}
	}()

	comment := &model.Comment{
		Content: content,
		PostID:  postId,
	}
	if err := s.commentDao.CreateComment(tx, comment); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) DeleteComment(cid uint) error {
	return s.commentDao.DeleteComment(cid)
}

