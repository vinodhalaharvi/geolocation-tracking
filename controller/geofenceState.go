package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vinodhalaharvi/geolocation-tracking/location"
	"github.com/vinodhalaharvi/geolocation-tracking/service"
	"net/http"
)

type GeoFenceStateController struct {
	service *service.GeoFenceStateService
}

func NewGeoFenceStateController(service *service.GeoFenceStateService) *GeoFenceStateController {
	return &GeoFenceStateController{service: service}
}

func (controller *GeoFenceStateController) CreatePolygon(c *gin.Context) {
	var polygon location.Polygon
	if err := c.ShouldBindJSON(&polygon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	controller.service.CreatePolygon(&polygon)

	c.JSON(http.StatusCreated, &polygon)
}

// ReadPolygon handles GET requests for a single polygon by ID
func (controller *GeoFenceStateController) ReadPolygon(c *gin.Context) {
	id := c.Param("id")
	polygon := controller.service.ReadPolygon(id)
	if polygon == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Polygon not found"})
		return
	}
	c.JSON(http.StatusOK, polygon)
}

// UpdatePolygon handles PUT requests to update a polygon
func (controller *GeoFenceStateController) UpdatePolygon(c *gin.Context) {
	id := c.Param("id")
	var updatedPolygon location.Polygon
	if err := c.ShouldBindJSON(&updatedPolygon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.service.UpdatePolygon(id, &updatedPolygon)
	c.JSON(http.StatusOK, gin.H{"message": "Polygon updated successfully"})
}

// DeletePolygon handles DELETE requests to remove a polygon
func (controller *GeoFenceStateController) DeletePolygon(c *gin.Context) {
	id := c.Param("id")
	controller.service.DeletePolygon(id)
	c.JSON(http.StatusOK, gin.H{"message": "Polygon deleted successfully"})
}

// ReadAsset Define other handlers for CRUD operations
func (controller *GeoFenceStateController) ReadAsset(c *gin.Context) {
	id := c.Param("id")
	car := controller.service.ReadAsset(id)
	if car == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	c.JSON(http.StatusOK, car)
}

func (controller *GeoFenceStateController) AddAsset(c *gin.Context) {
	var request struct {
		PolygonId string             `json:"polygonId"`
		Location  *location.Location `json:"location"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure Location is properly initialized; for example, set its ID or timestamp if necessary.
	request.Location.Id = uuid.New().String() // Optionally generate a unique ID for the location
	// Timestamp can be set here if relevant

	car := controller.service.AddAsset(request.PolygonId, request.Location)
	c.JSON(http.StatusOK, car)
}

func (controller *GeoFenceStateController) GetAllPolygons(context *gin.Context) {
	context.JSON(http.StatusOK, controller.service.GetPolygons())
}
