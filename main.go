package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/my_oauth_example/oauth"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// 读取环境变量
	githubClientId := os.Getenv("github_client_id")
	githubClientSecret := os.Getenv("github_client_secret")
	githubCallbackUrl := os.Getenv("github_callback_url")
	oauth.SetupGithubHandler(githubClientId, githubClientSecret, githubCallbackUrl)
	linuxDoClientId := os.Getenv("linuxdo_client_id")
	linuxDoClientSecret := os.Getenv("linuxdo_client_secret")
	linuxDoCallbackUrl := os.Getenv("linuxdo_callback_url")
	oauth.SetupLinuxDoHandler(linuxDoClientId, linuxDoClientSecret, linuxDoCallbackUrl)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))
	SetupRoutes(r)
	r.Run(":8080")
}
