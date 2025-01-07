package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	stdurl "net/url"
)

var LinuxDoOauthHandler LinuxDoHandler

func SetupLinuxDoHandler(clientId, clientSec, callback string) {
	proxyURL, err := stdurl.Parse("sock5://127.0.0.1:7890") // 替换为你的代理地址和端口
	if err != nil {
		panic(err)
	}
	LinuxDoOauthHandler = LinuxDoHandler{
		ClientId:  clientId,
		ClientSec: clientSec,
		Callback:  callback,
		oauth2Config: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSec,
			Scopes:       []string{"user:email", "user:username"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://connect.linux.do/oauth2/authorize",
				TokenURL: "https://connect.linux.do/oauth2/token",
			},
			RedirectURL: callback,
		},
		proxyUrl: *proxyURL,
	}
}

type LinuxDoHandler struct {
	Handler
	ClientId     string
	ClientSec    string
	Callback     string
	oauth2Config *oauth2.Config
	proxyUrl     stdurl.URL
}

func (h *LinuxDoHandler) Redirect2Oauth(ctx *gin.Context) {
	state, err := generateState()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state: " + err.Error()})
		return
	}
	session := sessions.Default(ctx)
	fmt.Println("state", state)
	session.Set("oauth_state", state)
	session.Save()
	url := h.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *LinuxDoHandler) GetCallback(ctx *gin.Context) {
	session := sessions.Default(ctx)
	code := ctx.Query("code")
	state := ctx.Query("state")
	oauthState := session.Get("oauth_state")
	fmt.Println("state2", state)
	// 校验 state
	if oauthState == nil || state != oauthState {
		ctx.String(http.StatusUnauthorized, "State value does not match")
		return
	}
	// 获取token
	// 使用code换取token
	//data := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", code, RedirectURI)
	//req, err := http.NewRequest("POST", TokenEndpoint, bytes.NewBuffer([]byte(data)))
	//if err != nil {
	//	return "", err
	//}
	//
	//req.SetBasicAuth(ClientID, ClientSecret)
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Accept", "application/json")
	token, err := h.oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code: " + err.Error()})
		return
	}
	client := h.oauth2Config.Client(ctx, token)
	//req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Get("https://connect.linux.do/api/user")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	var user map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
