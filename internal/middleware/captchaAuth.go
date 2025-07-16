package middleware

import (
	"CloudPhoto/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func X() {}
func CaptchaAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		print("captcha auth")
		captchaId := c.GetHeader("captcha-id")
		captchaCode := c.GetHeader("captcha-code")

		//if storage.GetCaptcha(captchaId) != captchaCode {
		//	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		//	c.Abort()
		//	return
		//}
		//storage.RemoveCaptcha(captchaId)
		if (*storage.GetCaptcha()).Verify(captchaId, captchaCode, true) {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
		}
	}
}
