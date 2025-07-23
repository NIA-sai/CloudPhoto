package server

import (
	"CloudPhoto/cmd/daemon/cleaner"
	"CloudPhoto/config"
	"CloudPhoto/internal/database"
	"CloudPhoto/internal/middleware"
	"CloudPhoto/internal/module"
	"CloudPhoto/internal/storage"
	"CloudPhoto/internal/tool"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func initialize() {
	config.Read()
	database.Init()
	var f = []string{
		config.Get().App.StaticRoot + storage.CutOutFilePath,
	}
	cleaner.StartFileCleanupTask(f, 24*time.Hour)
}

var r = gin.New()
var srv *http.Server

func Start() {
	//初始化基本模块
	initialize()
	//添加中间件
	r.Use(
		middleware.Cors,
	)
	r.Static(config.Get().App.StaticRelativePath, config.Get().App.StaticRoot)
	for _, m := range module.Modules {
		m.Init()
		m.InitRouter(r.Group(m.GetName()))
	}
	srv = &http.Server{
		Addr:         config.Get().App.Host + ":" + strconv.Itoa(config.Get().App.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	tool.PanicIfErr(srv.ListenAndServe())
}
func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
