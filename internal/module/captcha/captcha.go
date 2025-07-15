package captcha

import (
	"CloudPhoto/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

// imageCaptcha base64:<img :src="'data:image/png;base64,' + captchaImage" />
func imageCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverString(80, 200, 5, 3, 4,
		"1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		nil, nil, []string{})
	captcha := base64Captcha.NewCaptcha(driver, *storage.GetCaptcha())
	id, b64s, _, err := captcha.Generate()

	if err != nil {
		print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "captcha generate failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
}
