package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/vinodhalaharvi/geolocation-tracking/controller"
	"github.com/vinodhalaharvi/geolocation-tracking/service"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("location-tracking", store))

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Routes
	r.GET("/login", controller.LoginHandler)
	r.POST("/login", controller.SimpleLoginHandler)
	r.GET("/auth/google/login", controller.GoogleLoginHandler)
	r.GET("/auth/google/callback", controller.GoogleAuthCallback)

	authorized := r.Group("/")
	authorized.Use(controller.GoogleOAuthMiddleware())

	authorized.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	logout := controller.Logout{}

	r.GET("/logout", logout.Handle)

	r.GET("/error", func(c *gin.Context) {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
	})

	authorized.GET("/js", func(c *gin.Context) {
		c.HTML(http.StatusOK, "js.html", nil)
	})

	authorized.GET("/maps", func(c *gin.Context) {
		c.HTML(http.StatusOK, "maps.html", nil)
	})

	handler := controller.NewWebsocket()
	r.GET("/ws", func(context *gin.Context) {
		// add to the clients
		handler.Handle(context.Writer, context.Request)
	})

	fmt.Printf("http://localhost:8080\n")

	geoFenceStateService := service.NewGeoFenceStateService()
	geoFenceStateController := controller.NewGeoFenceStateController(geoFenceStateService)

	// Polygon routes
	authorized.GET("/polygons", geoFenceStateController.GetAllPolygons)
	authorized.POST("/polygons", geoFenceStateController.CreatePolygon)
	authorized.POST("/addAsset", geoFenceStateController.AddAsset)
	authorized.GET("/polygons/:id", geoFenceStateController.ReadPolygon)
	authorized.PUT("/polygons/:id", geoFenceStateController.UpdatePolygon)
	authorized.DELETE("/polygons/:id", geoFenceStateController.DeletePolygon)

	// Asset routes
	authorized.GET("/cars/:id", geoFenceStateController.ReadAsset)

	authorized.POST("/simulate", func(c *gin.Context) {
		// Start simulation logic here
		go handler.RunSimulation(c, geoFenceStateService)
		c.JSON(http.StatusOK, gin.H{"message": "Simulation started"})
	})

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
