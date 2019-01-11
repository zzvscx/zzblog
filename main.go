package main

import (
	"fmt"
	_ "github.com/gin-gonic/gin"
	"net/http"
	"zzblog/models"
	"zzblog/pkg/setting"
	"zzblog/routers"
)

func main() {
	setting.Setup()
	models.Setup()
	defer models.CloseDB()

	routersInit := routers.InitRouter()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", setting.HTTPPort),
		Handler: routersInit,
		//ReadTimeout:    setting.ReadTimeout,
		//WriteTimeout:   setting.WriteTimeout,
		//MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
