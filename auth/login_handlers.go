package auth

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"html/template"
	"io"
	"net/http"
)

func homeHandler(c *gin.Context) {
	session := sessions.Default(c)
	userEmail := session.Get("user-email")

	if userEmail != nil {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"email": userEmail})
	} else {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	}
}

func LoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func GoogleLoginHandler(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleAuthCallback(c *gin.Context) {
	session := sessions.Default(c)

	// Handle the exchange code to initiate a transport.
	token, err := googleOauthConfig.Exchange(c, c.Query("code"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(contents, &userInfo); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Store the user's email in session
	session.Set("user-email", userInfo.Email)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func HtmlTemplate() *template.Template {
	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}
	return t
}

func GoogleOAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user-email")

		// If the user is not logged in, redirect to login
		if user == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}

		// Otherwise, proceed to the next middleware
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user-email")

		// If the user is not logged in, redirect to login
		if user == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}

		// Otherwise, proceed to the next middleware
		c.Next()
	}
}
