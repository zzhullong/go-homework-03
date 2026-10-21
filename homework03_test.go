package main

import (
	"homework03/internal/model"
	"testing"

	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func getTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:testdb?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}
	err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err != nil {
		t.Fatalf("迁移建表失败: %v", err)
	}
	return db
}

func TestModelDefine(t *testing.T) {
	db := getTestDB(t)
	user := model.User{Name: "testuser", Pwd: "123456"}
	err := db.Create(&user).Error
	if err != nil {
		t.Fatalf("创建用户失败：%v", err)
	}
	var u model.User
	db.First(&u, user.ID)
	if u.Name != "testuser" {
		t.Error("用户读取数据不匹配")
	}
	t.Log("✅ 题目1：模型建表测试通过")
}

func TestRelateQuery(t *testing.T) {
	db := getTestDB(t)
	user := model.User{Name: "alice", Pwd: "111"}
	db.Create(&user)

	p1 := model.Post{UserID: user.ID, Title: "文章A", Content: "内容A", CommentCount: 2, CommentStatus: "有评论"}
	p2 := model.Post{UserID: user.ID, Title: "文章B", Content: "内容B", CommentCount: 0, CommentStatus: "无评论"}
	db.Create(&p1)
	db.Create(&p2)

	c1 := model.Comment{PostID: p1.ID, Content: "很好看"}
	c2 := model.Comment{PostID: p1.ID, Content: "不错"}
	db.Create(&c1)
	db.Create(&c2)

	var u model.User
	err := db.Preload("Posts.Comments").First(&u, user.ID).Error
	if err != nil {
		t.Fatal(err)
	}
	if len(u.Posts) != 2 {
		t.Errorf("预期2篇文章，实际:%d", len(u.Posts))
	}
	if len(u.Posts[0].Comments) != 2 {
		t.Errorf("文章A预期2条评论，实际:%d", len(u.Posts[0].Comments))
	}

	var maxPost model.Post
	err = db.Order("comment_count desc").First(&maxPost).Error
	if err != nil {
		t.Fatal(err)
	}
	if maxPost.Title != "文章A" {
		t.Error("评论最多文章查询错误")
	}
	t.Log("✅ 题目2：关联查询测试通过")
}
func TestHook(t *testing.T) {
	db := getTestDB(t)
	user := model.User{Name: "hooktest", Pwd: "pwd"}
	db.Create(&user)
	post := model.Post{UserID: user.ID, Title: "钩子测试文章", Content: "test content"}
	err := db.Create(&post).Error
	if err != nil {
		t.Fatal(err)
	}

	var u model.User
	err = db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&u, user.ID).Error
	if err != nil {
		t.Fatal(err)
	}
	if u.Count != 1 {
		t.Errorf("钩子失败，用户文章数预期1，实际=%d", u.Count)
	}

	c1 := model.Comment{PostID: post.ID, Content: "评论1"}
	c2 := model.Comment{PostID: post.ID, Content: "评论2"}
	db.Create(&c1)
	db.Create(&c2)

	// ========= 删除第一条评论 c1 =========
	db.Delete(&c1)
	var p1 model.Post
	err = db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&p1, post.ID).Error
	if err != nil {
		t.Fatal(err)
	}
	if p1.CommentCount != 1 || p1.CommentStatus != "有评论" {
		t.Errorf("删除一条评论，状态或数量错误, 预期 count=1,status=有评论，实际 count=%d, status=%s", p1.CommentCount, p1.CommentStatus)
	}

	// ========= 删除第二条评论 c2 =========
	db.Delete(&c2)
	var p2 model.Post
	err = db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&p2, post.ID).Error
	if err != nil {
		t.Fatal(err)
	}
	if p2.CommentCount != 0 || p2.CommentStatus != "无评论" {
		t.Errorf("全部评论删除后，状态没有变成无评论, 预期 count=0,status=无评论，实际 count=%d, status=%s", p2.CommentCount, p2.CommentStatus)
	}

	t.Log("✅ 题目3：钩子函数测试通过")
}