package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zzblog/models"
)

func PostGet(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.HTML(http.StatusOK, "post/new.html", gin.H{
		"post": post,
	})
}

func PostNew(c *gin.Context) {
	c.HTML(http.StatusOK, "post/new.html", nil)
}

func PostCreate(c *gin.Context) {
	title := c.PostForm("title")
	body := c.PostForm("body")
	post := &models.Post{
		Title: title,
		Body:  body,
	}
	err := post.Insert()
	if err != nil {
		c.HTML(http.StatusOK, "post/new.html", gin.H{
			"post":    post,
			"message": err.Error(),
		})
	}
	c.Redirect(http.StatusMovedPermanently, "/post/"+strconv.FormatUint(uint64(post.ID), 10))
}

func PostEdit(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.HTML(http.StatusOK, "post/modify.html", gin.H{
		"post": post,
	})
}

func PostUpdate(c *gin.Context) {
	id := c.Param("id")
	title := c.PostForm("title")
	body := c.PostForm("body")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	post := &models.Post{
		Title: title,
		Body:  body,
	}
	post.ID = uint(pid)
	err = post.Update()
	if err != nil {
		c.HTML(http.StatusOK, "post/modify.html", gin.H{
			"post":    post,
			"message": err.Error(),
		})
	}
	c.Redirect(http.StatusMovedPermanently, "/post/"+id)
}

func PostDelete(c *gin.Context) {
	id := c.PostForm("id")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	post := &models.Post{}
	post.ID = uint(pid)
	err = post.Delete()
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/post")
}

func PostIndex(c *gin.Context) {
	posts, err := models.ListPostByTagId("")
	c.HTML(http.StatusOK, "", gin.H{
		"posts":   posts,
		"message": err.Error(),
	})
}
