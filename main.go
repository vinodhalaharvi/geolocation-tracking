package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/vinodhalaharvi/geolocation-tracking/auth"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("location-tracking", store))

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Routes
	r.GET("/login", auth.LoginHandler)
	r.GET("/auth/google/login", auth.GoogleLoginHandler)
	r.GET("/auth/google/callback", auth.GoogleAuthCallback)

	authorized := r.Group("/")
	authorized.Use(auth.GoogleOAuthMiddleware())

	authorized.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		err := session.Save()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
		}
		c.Redirect(http.StatusSeeOther, "/login")
	})

	r.GET("/error", func(c *gin.Context) {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
	})

	authorized.GET("/maps", func(c *gin.Context) {
		c.HTML(http.StatusOK, "maps.html", nil)
	})

	handler := WsHandler{}
	r.GET("/ws", func(context *gin.Context) {
		handler.Handle(context.Writer, context.Request)
	})

	fmt.Printf("http://localhost:8080\n")

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
