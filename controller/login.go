package controller

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
)

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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(response.Body)

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
	err = session.Save()
	if err != nil {
		log.Printf("Failed to save session: %v", err)
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func SimpleLoginHandler(c *gin.Context) {
	session := sessions.Default(c)

	// Normally, you'd get these from the request, e.g., c.PostForm("username")
	// Here, we're hardcoding them for demonstration purposes
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Placeholder for proper credential validation
	// In a real app, replace this with database lookup and secure password comparison
	const expectedUsername = "admin"
	const expectedPassword = "adminPass"

	// Validate credentials
	if username == expectedUsername && password == expectedPassword {
		// Authentication successful
		session.Set("isAuthenticated", true)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		c.Redirect(http.StatusFound, "/")
	} else {
		// Authentication failed
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
	}
}

func GoogleOAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user-email")
		isAuthenticated := session.Get("isAuthenticated") // This could be set by your SimpleLoginHandler

		// Check if the user is logged in through any method
		if user == nil && isAuthenticated != true {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}

		// Otherwise, proceed to the next middleware
		c.Next()
	}
}
