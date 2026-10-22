package dao

import (
	"homework03/internal/model"
	"homework03/pkg/db"

	"gorm.io/gorm"
)

type CommentDao struct{}

func NewCommentDao() *CommentDao {
	return &CommentDao{}
}

// CreateComment 新增评论
func (d *CommentDao) CreateComment(tx *gorm.DB,comment *model.Comment) error {
	return tx.Create(comment).Error
}

// DeleteComment 删除评论
func (d *CommentDao) DeleteComment(commentId uint) error {
	return db.DB.Delete(&model.Comment{}, commentId).Error
}

