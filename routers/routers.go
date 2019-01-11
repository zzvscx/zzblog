package routers

import (
	"github.com/gin-gonic/gin"
	"zzblog/pkg/setting"
	"zzblog/views"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(setting.RunMode)
	post := r.Group("/post")
	{
		post.GET("/", views.PostIndex)
		post.GET("/new", views.PostNew)
		post.POST("/new", views.PostCreate)
		post.GET("/edit/:id", views.PostEdit)
		post.POST("/edit/:id", views.PostUpdate)
		post.POST("/delete/:id", views.PostDelete)
	}

	tag := r.Group("/tag")
	{
		tag.POST("/new", views.TagCreate)
	}
	return r
}
