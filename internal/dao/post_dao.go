package dao

import (
	"errors"
	"homework03/internal/model"
	"homework03/pkg/db"

	"gorm.io/gorm"
)

type PostDao struct{
	db *gorm.DB
}

func NewPostDao() *PostDao {
	return &PostDao{}
}

func (d *PostDao) GetDB() *gorm.DB {
	return d.db
}

// GetMostCommentPost 获取评论数量最多的文章
func (d *PostDao) GetMostCommentPost() (*model.Post, error) {
	var post model.Post
	err := db.DB.Order("comment_count desc, id").First(&post).Error
	if err != nil {
		// 如果只是找不到记录，返回 nil,nil，代表暂无数据
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// 其余数据库错误正常返回
		return nil, err
	}
	return &post, nil
}

//  用户文章计数原子+N
func (d *PostDao) IncArticleCount(tx *gorm.DB, postId int, num int) error {
	return tx.Model(&model.Post{}).
		Where("id = ?", postId).
		UpdateColumn("CommentCount", gorm.Expr("CommentCount + ?", num)).
		Error
}


// CreatePost 创建文章
func (d *PostDao) CreatePost(tx *gorm.DB, post *model.Post) error {
	return tx.Create(post).Error
}

// GetPostByID 根据id获取文章
func (d *PostDao) GetPostByID(postId int) (*model.Post, error) {
	var post model.Post
	err := db.DB.First(&post, postId).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// DeletePost 软删除文章
func (d *PostDao) DeletePost(id int) error {
	return db.DB.Delete(&model.Post{}, id).Error
}