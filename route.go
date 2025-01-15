package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/my_oauth_example/oauth"
	"net/http"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", gin.H{
			"title": "Home",
		})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// github
	r.GET("/login/github", oauth.GithubOauthHandler.Redirect2Oauth)
	r.GET("v1/auth/oauth2/github/callback", oauth.GithubOauthHandler.GetCallback)
	// linux do
	r.GET("/login/linuxdo", oauth.LinuxDoOauthHandler.Redirect2Oauth)
	r.GET("/v1/auth/oauth2/linux_do/callback", oauth.LinuxDoOauthHandler.GetCallback)
	// ...
}
