package middleware

import (
	"github.com/gin-contrib/cors"
	"time"
)

var Cors = cors.New(cors.Config{
	AllowAllOrigins: true,
	AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:    []string{"Origin", "Content-Type", "captcha-id", "captcha-code"},
	MaxAge:          24 * time.Hour,
})
