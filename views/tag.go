package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zzblog/models"
)

func TagCreate(c *gin.Context) {
	name := c.PostForm("name")
	tag := &models.Tag{Name: name}
	err := tag.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data": tag,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	}
}

func TagGet(c *gin.Context) {
	id := c.Param("id")
	posts, err := models.ListPostByTagId(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"posts": posts,
	})
}
