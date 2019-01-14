package views

import "C"
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zzblog/models"
)

func PageGet(c *gin.Context) {
	id := c.Param("id")
	page, err := models.GetPageById(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.HTML(http.StatusOK, "page/display.html", gin.H{
		"page": page,
	})
}

func PageNew(c *gin.Context) {
	c.HTML(http.StatusOK, "page/new.html", nil)
}

func PageCreate(c *gin.Context) {
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	page := &models.Page{
		Body:        body,
		IsPublished: isPublished,
	}
	err := page.Insert()
	if err != nil {
		c.HTML(http.StatusOK, "page/new.html", gin.H{
			"message": err.Error(),
			"page":    page,
		})
	}
	c.Redirect(http.StatusMovedPermanently, "/page/"+strconv.FormatUint(uint64(page.ID), 10))
}

func PageEdit(c *gin.Context) {
	id := c.Param("id")
	page, err := models.GetPageById(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.HTML(http.StatusOK, "page/modufy.html", gin.H{
		"page": page,
	})
}

func PageUpdate(c *gin.Context) {
	id := c.Param("id")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	page := &models.Page{Body: body, IsPublished: isPublished}
	page.ID = uint(pid)
	err = page.Update()
	if err != nil {
		//TODO
	}
	c.Redirect(http.StatusMovedPermanently, "/page/"+id)
}

func PageDelete(c *gin.Context) {
	id := c.PostForm("id")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	page := &models.Page{}
	page.ID = uint(pid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Redirect(http.StatusMovedPermanently, "/page/")
}

func PageIndex(c *gin.Context) {
	pages, _ := models.ListPage()
	c.HTML(http.StatusOK, "page/", gin.H{
		"pages": pages,
	})
}
