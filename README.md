# go-homework-03 

进阶gorm


项目结构：
env 环境变量
config
    基础相关配置，包括配置和读取
internal 
    dao，model，service 分层设计
pkg/db
    全局DB初始化

注意事项：
1.环境适配
windows/macOS 支持的db驱动不一样， windows支持github.com/glebarez/go-sqlite  macOS和linux支持gorm.io/driver/sqlite

如要适配当前操作系统环境， 有几个文件需要修改一下：

go.mod require地址修改对应地址
homework03_test.go import修改对应包
pkg/db/db.go   import修改对应包



验证方式:
1.homework03_main.go 创建数据，查询数据
2.homework03_test.go  test题目123通过性


题目1：模型定义
假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
要求 ：
使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
编写Go代码，使用Gorm创建这些模型对应的数据库表。



题目2：关联查询
基于上述博客系统的模型定义。
要求 ：
编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
编写Go代码，使用Gorm查询评论数量最多的文章信息。



题目3：钩子函数
继续使用博客系统的模型。
要求 ：
为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

