package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// 配置结构体
type Config struct {
	gorm.Model
	SearchEngine string
	Footer       string
}

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("config.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&Config{})
}

func main() {
	r := gin.Default()

	// 首页处理
	r.GET("/", func(c *gin.Context) {
		var config Config
		db.First(&config)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Webstack-Go 星空主题",
			"footer":     config.Footer,
			"search_eng": config.SearchEngine,
		})
	})

	// 后台管理页面
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.html", nil)
	})

	// 保存后台配置
	r.POST("/admin/save-config", func(c *gin.Context) {
		searchEngine := c.PostForm("search_engine")
		footer := c.PostForm("footer")

		config := Config{SearchEngine: searchEngine, Footer: footer}
		db.Create(&config)

		c.Redirect(http.StatusFound, "/admin")
	})

	// 启动服务器
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.Run(":8080")
}
