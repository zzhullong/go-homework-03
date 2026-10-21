package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	ID    int   `gorm:"primaryKey"`
	Name  string `gorm:"size:64;not null"`
	Email *string `gorm:"size:128"`
	Pwd string   `gorm:"size:64;not null"`
	Posts []Post
	Count int  `gorm:"column:count;default:0"`
}

type Post struct {
	gorm.Model
	ID            int   `gorm:"primaryKey"`
	Title         string `gorm:"size:100;not null"`
	Content       string `gorm:"type:text;not null"`
	UserID        int   `gorm:"not null"`
	CommentCount  int    `gorm:"default:0"`
	CommentStatus string `gorm:"size:20;default:'无评论'"`
	Comments      []Comment
}

type Comment struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"size:500;not null"`
	PostID  int   `gorm:"not null"`
}



func (p *Post) AfterCreate(tx *gorm.DB) error {
	var u User
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&u, p.UserID).Error
	if err != nil {
		return err
	}
	u.Count++
	return tx.Save(&u).Error
}

// 新增评论：AfterCreate，给帖子数量 +1
func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	var post Post
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&post, c.PostID).Error
	if err != nil {
		return err
	}
	post.CommentCount++
	post.CommentStatus = "有评论"
	return tx.Save(&post).Error
}

// 删除评论：BeforeDelete，给帖子数量 -1（软删除会触发 BeforeDelete）
func (c *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	var post Post
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&post, c.PostID).Error
	if err != nil {
		return err
	}
	// 防止负数
	if post.CommentCount > 0 {
		post.CommentCount--
	}
	if post.CommentCount == 0 {
		post.CommentStatus = "无评论"
	} else {
		post.CommentStatus = "有评论"
	}
	return tx.Save(&post).Error
}