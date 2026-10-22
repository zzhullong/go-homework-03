package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	
	Name  string `gorm:"size:64;not null"`
	Email *string `gorm:"size:128"`
	Pwd string   `gorm:"size:64;not null"`
	Posts []Post
	Count int  `gorm:"column:count;default:0"`
}

type Post struct {
	gorm.Model

	Title         string `gorm:"size:100;not null"`
	Content       string `gorm:"type:text;not null"`
	UserID        uint   `gorm:"not null"`
	CommentCount  int    `gorm:"default:0"`
	CommentStatus string `gorm:"size:20;default:'无评论'"`
	Comments      []Comment
}

type Comment struct {
	gorm.Model
	
	Content string `gorm:"size:500;not null"`
	PostID  uint   `gorm:"not null"`
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

// 删除评论：BeforeDelete，给帖子数量 -1（软删除会触发 BeforeDelete），假设钩子一次只能删除一条的场景
func (c *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	var post Post

	
	// 判断：当前是不是【单实例删除】，不是就跳过业务逻辑，放行删除
	//假设postid 存在0的情况
	//
	// val := tx.Statement.ReflectValue
	// isBatchTemplate := c.PostID == 0 && val.IsValid() && val.Kind() == reflect.Struct && val.IsZero()
	isBatchTemplate := c.ID==0 

	if isBatchTemplate {
		return nil
	}

	


	//增加判断是否已经软删除了，如果软删除的话，不再触发钩子进行-1,默认软删除就会有deleteat时间
	if post.DeletedAt.Valid {
		return nil
	}

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