package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vinodhalaharvi/geolocation-tracking/location"
	"github.com/vinodhalaharvi/geolocation-tracking/service"
	"net/http"
)

type AssetController struct {
	assetService *service.AssetService
}

func NewAssetController(cs *service.AssetService) *AssetController {
	return &AssetController{
		assetService: cs,
	}
}

func (cc *AssetController) CreateAsset(c *gin.Context) {
	var car location.Asset
	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cc.assetService.AddAsset(&car)

	c.JSON(http.StatusCreated, &car)
}

// ReadAsset retrieves a car by its ID
func (cc *AssetController) ReadAsset(c *gin.Context) {
	id := c.Param("id")
	car, found := cc.assetService.GetAsset(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	c.JSON(http.StatusOK, car)
}

// UpdateAsset updates a car by its ID
func (cc *AssetController) UpdateAsset(c *gin.Context) {
	var newAsset location.Asset
	if err := c.ShouldBindJSON(&newAsset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	updated := cc.assetService.UpdateAsset(id, &newAsset)
	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	c.JSON(http.StatusOK, &newAsset)
}

// DeleteAsset deletes a car by its ID
func (cc *AssetController) DeleteAsset(c *gin.Context) {
	id := c.Param("id")
	deleted := cc.assetService.DeleteAsset(id)
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Asset deleted"})
}
