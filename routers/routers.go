package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"zzblog/pkg/setting"
	"zzblog/pkg/utils"
	"zzblog/views"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	setTemplate(r)
	setSessions(r)
	r.Static("/static", "./static")
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(setting.RunMode)
	r.GET("/", views.IndexGet)
	post := r.Group("/post")
	{
		post.GET("/", views.PostIndex)
		post.GET("/new", views.PostNew)
		post.POST("/new", views.PostCreate)
		post.GET("/edit/:id", views.PostEdit)
		post.POST("/edit/:id", views.PostUpdate)
		post.POST("/delete/:id", views.PostDelete)
	}
	page := r.Group("/page")
	{
		page.GET("/", views.PageIndex)
		page.GET("/new", views.PageNew)
		page.POST("/new", views.PageCreate)
		page.GET("/edit/:id", views.PageEdit)
		page.POST("/edit/:id", views.PageUpdate)
		page.POST("/delete/:id", views.PageDelete)
	}

	tag := r.Group("/tag")
	{
		tag.POST("/new", views.TagCreate)
		tag.GET("/:id", views.TagGet)
	}
	archives := r.Group("/archives")
	{
		archives.GET("/:year/:month", views.ArchiveGet)
	}
	user := r.Group("/user")
	{
		user.GET("/signin", views.SigninGet)
		user.POST("/signin", views.SigninPost)
		user.GET("/signup", views.SignupGet)
		user.POST("/signup", views.SignupPost)
		user.GET("/logout", views.LogoutGet)
	}

	return r
}

func setTemplate(r *gin.Engine) {
	funcMap := template.FuncMap{
		"dateFormat": utils.DateFormat,
		"substring":  utils.Substring,
		"isOdd":      utils.IsOdd,
		"isEven":     utils.IsEven,
	}
	r.FuncMap = funcMap
	r.LoadHTMLGlob("template/**/*")
}

func setSessions(r *gin.Engine) {
	store := cookie.NewStore([]byte(setting.SessionSecret))
	r.Use(sessions.Sessions("mysession", store))
}
