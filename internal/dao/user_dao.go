package dao

import (
	"homework03/internal/model"
	"homework03/pkg/db"
)

type UserDao struct{}

func NewUserDao() *UserDao {
	return &UserDao{}
}

// GetUserWithPostsAndComments 根据用户ID 查询用户，预加载文章和文章评论（一对多）
func (d *UserDao) GetUserWithPostsAndComments(uid uint) (*model.User, error) {
	var user model.User
	// Preload 嵌套预加载 Posts.Comments
	err := db.DB.Preload("Posts.Comments").First(&user, uid).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 新建用户
func (d *UserDao) CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}

func (d *UserDao) DeleteUser(id uint) error {
	return db.DB.Delete(&model.User{}, id).Error
}


// //  用户文章计数原子+N
// func (d *UserDao) IncArticleCount(tx *gorm.DB, userId int, num int) error {
// 	return tx.Model(&model.User{}).
// 		Where("id = ?", userId).
// 		UpdateColumn("count", gorm.Expr("count + ?", num)).
// 		Error
// }


