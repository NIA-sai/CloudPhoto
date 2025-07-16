package server

import (
	"CloudPhoto/config"
	"CloudPhoto/internal/database"
	"CloudPhoto/internal/module"
	"CloudPhoto/internal/tool"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func initialize() {
	config.Read()
	database.Init()
}

var r = gin.New()
var srv *http.Server

func Start() {
	initialize()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "captcha-id", "captcha-code"},
		MaxAge:          24 * time.Hour,
	}))
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
