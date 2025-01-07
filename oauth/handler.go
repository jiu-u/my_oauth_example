package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
}

func (h *Handler) GotoLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "home.html", gin.H{
		"title": "Home",
	})
}

// generateState 生成一个随机的state参数
func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
