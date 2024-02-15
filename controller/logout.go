package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Logout struct {
}

func (l *Logout) Handle(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		// Consider logging the error here as well
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return // Ensure you return after sending the response
	}

	// Prevent browsers from caching the response to this redirect
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// Redirect to Google logout
	//c.Redirect(http.StatusTemporaryRedirect, "https://www.google.com/accounts/Logout")
	c.Redirect(http.StatusSeeOther, "/login")
}
