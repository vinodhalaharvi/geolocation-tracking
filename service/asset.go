// service/carService.go

package service

import (
	"github.com/vinodhalaharvi/geolocation-tracking/location"
	"sync"
)

type AssetService struct {
	assets []*location.Asset
	mu     sync.Mutex
}

func NewAssetService() *AssetService {
	return &AssetService{
		assets: []*location.Asset{},
	}
}

func (cs *AssetService) AddAsset(car *location.Asset) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.assets = append(cs.assets, car)
}

func (cs *AssetService) GetAsset(id string) (*location.Asset, bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	for _, car := range cs.assets {
		if car.Id == id {
			return car, true
		}
	}

	return &location.Asset{}, false
}

func (cs *AssetService) UpdateAsset(id string, newAsset *location.Asset) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	for i, car := range cs.assets {
		if car.Id == id {
			cs.assets[i] = newAsset
			return true
		}
	}

	return false
}

func (cs *AssetService) DeleteAsset(id string) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	for i, car := range cs.assets {
		if car.Id == id {
			cs.assets = append(cs.assets[:i], cs.assets[i+1:]...)
			return true
		}
	}

	return false
}

func (cs *AssetService) GetAllAssets() []*location.Asset {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	return cs.assets
}
