package core

import (
	"sync"

	"github.com/pakerfeldt/lovi/pkg/models"
)

type UnAckedAlertsMap struct {
	sync.RWMutex
	internal map[string]models.Alert
}

func NewUnAckedAlertsMap() *UnAckedAlertsMap {
	return &UnAckedAlertsMap{
		internal: make(map[string]models.Alert),
	}
}

func (rm *UnAckedAlertsMap) Load(key string) (value models.Alert, ok bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}

func (rm *UnAckedAlertsMap) Delete(key string) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *UnAckedAlertsMap) Store(key string, value models.Alert) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}
