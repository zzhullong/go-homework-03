package main

import (
	"homework03/config"
	"homework03/internal/service"
	"homework03/pkg/db"
	"log"
	"strconv"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("加载配置失败：", err)
	}
	log.Printf("✅ driver=%s, DSN=%s", cfg.DB.Driver, cfg.DB.DSN)

	err = db.Init(db.Config{
		Driver: cfg.DB.Driver,
		DSN:    cfg.DB.DSN,
	})
	if err != nil {
		log.Fatal("数据库初始化失败：", err)
	}
	log.Println("数据库初始化成功")

	postSvc := service.NewPostService()
	userSvc := service.NewUserService()
	commentSvc := service.NewCommentService()


	//postSvc.DeletePost(1)

	//创建用户及关联数据
	user, err := userSvc.CreateUser("zhangsan", "123456")
	if err != nil {
		log.Fatalf("创建用户失败：%v", err)
	}
	log.Printf("✅ 创建用户成功，用户ID：%d", user.ID)

	
for i := 1; i < 11; i++ {
	title := "我的第" + strconv.Itoa(i) + "篇Go博客"
	content := "GORM一对多学习内容"

	post, err := postSvc.CreatePost(title, content, user.ID)
	if err != nil {
		log.Printf("创建文章失败：%v", err)
		break
	}
	log.Printf("✅ 创建文章成功，文章ID：%d，所属用户ID：%d", post.ID, post.UserID)


}

for i := 1; i < 10; i++ {
		comment1, err := commentSvc.CreateComment("写得很棒，学到了！",i)
	if err != nil {
		log.Printf("创建评论失败：%v", err)
		break
	}
	log.Printf("✅ 创建评论1成功，评论ID：%d，文章ID：%d", comment1.ID, comment1.PostID)
}

	log.Printf("db全局实例地址：%p", db.DB) 
	
	userInfo, err := userSvc.GetUserWithPostsAndComments(1)
	if err != nil {
		log.Println("查询用户信息失败：", err)
	} else {
		log.Printf("用户信息：%+v", userInfo)
		log.Printf("用户文章总数：%+v", userInfo.Count)
		log.Printf("用户文章：%+v", userInfo.Posts)
	}

	hotPost, err := postSvc.GetMostCommentPost()
	if err != nil {
		log.Printf("查询评论最多文章，数据库出错：%v", err)
		return
	}
	if hotPost == nil {
		log.Println("暂无文章")
		return
	}
	log.Printf("🔥 评论最多文章：id=%d，标题=%s，评论数=%d", hotPost.ID, hotPost.Title, hotPost.CommentCount)
}