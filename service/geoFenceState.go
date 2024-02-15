package service

import (
	"github.com/google/uuid"
	"github.com/vinodhalaharvi/geolocation-tracking/location"
	"math/rand"
	"sync"
)

type GeoFenceStateService struct {
	State location.GeoFenceState `json:"state"`
	mu    sync.Mutex
}

func NewGeoFenceStateService() *GeoFenceStateService {
	return &GeoFenceStateService{
		State: location.GeoFenceState{
			Assets:   []*location.Asset{},
			Polygons: []*location.Polygon{},
		},
	}
}

func (gfs *GeoFenceStateService) ReadPolygon(id string) *location.Polygon {
	gfs.mu.Lock()
	defer gfs.mu.Unlock()

	for _, polygon := range gfs.State.Polygons {
		if polygon.Id == id {
			return polygon
		}
	}
	return nil
}

func (gfs *GeoFenceStateService) UpdatePolygon(id string, updatedPolygon *location.Polygon) {
	gfs.mu.Lock()
	defer gfs.mu.Unlock()

	for i, polygon := range gfs.State.Polygons {
		if polygon.Id == id {
			gfs.State.Polygons[i] = updatedPolygon
			break
		}
	}
}

func (gfs *GeoFenceStateService) DeletePolygon(id string) {
	gfs.mu.Lock()
	defer gfs.mu.Unlock()

	// Remove the polygon
	for i, polygon := range gfs.State.Polygons {
		if polygon.Id == id {
			gfs.State.Polygons = append(gfs.State.Polygons[:i], gfs.State.Polygons[i+1:]...)
			break
		}
	}

	// Optionally, remove assets associated with this polygon
	var newAssets []*location.Asset
	for _, car := range gfs.State.Assets {
		if car.PolygonId != id {
			newAssets = append(newAssets, car)
		}
	}
	gfs.State.Assets = newAssets
}

func (gfs *GeoFenceStateService) CreateAsset(car *location.Asset) {
	gfs.mu.Lock()
	defer gfs.mu.Unlock()

	car.Id = uuid.New().String() // Ensure the car has a unique ID
	gfs.State.Assets = append(gfs.State.Assets, car)
}

func (gfs *GeoFenceStateService) UpdateAsset(id string, updatedAsset *location.Asset) {
	gfs.mu.Lock()
	defer gfs.mu.Unlock()

	for i, car := range gfs.State.Assets {
		if car.Id == id {
			gfs.State.Assets[i] = updatedAsset
			break
		}
	}
}

func (gfs *GeoFenceStateService) ReadAsset(id string) *location.Asset {
	gfs.mu.Lock()
	defer gfs.mu.Unlock()

	for _, car := range gfs.State.Assets {
		if car.Id == id {
			return car
		}
	}
	return nil

}

func (gfs *GeoFenceStateService) CreatePolygon(l *location.Polygon) {
	gfs.mu.Lock() // Lock the mutex before modifying the state
	defer gfs.mu.Unlock()

	// Generate a unique ID for the new polygon
	l.Id = uuid.New().String()

	// Append the new polygon to the GeoFenceState
	gfs.State.Polygons = append(gfs.State.Polygons, l)
}

func (gfs *GeoFenceStateService) SimulateAssetMovements() {
	gfs.mu.Lock()
	defer gfs.mu.Unlock()

	for _, car := range gfs.State.Assets {
		// Simulate movement by adjusting the car's location slightly
		car.Location.Latitude += (rand.Float64() - 0.5) * 0.001 // Adjust these values based on your needs
		car.Location.Longitude += (rand.Float64() - 0.5) * 0.001
	}
}

func (gfs *GeoFenceStateService) AddAsset(polygonId string, loc *location.Location) *location.Asset {
	gfs.mu.Lock()
	defer gfs.mu.Unlock()

	car := &location.Asset{
		Id:        uuid.New().String(),
		Location:  loc, // loc is already a pointer to location.Location
		PolygonId: polygonId,
	}
	gfs.State.Assets = append(gfs.State.Assets, car)
	return car
}

func (gfs *GeoFenceStateService) GetPolygons() any {
	gfs.mu.Lock()
	defer gfs.mu.Unlock()
	return gfs.State.Polygons
}
