package core

import (
	"sync"

	"github.com/pakerfeldt/lovi/pkg/models"
)

type UnAckedEventsMap struct {
	sync.RWMutex
	internal map[string]models.Event
}

func NewUnAckedEventsMap() *UnAckedEventsMap {
	return &UnAckedEventsMap{
		internal: make(map[string]models.Event),
	}
}

func (rm *UnAckedEventsMap) Load(key string) (value models.Event, ok bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}

func (rm *UnAckedEventsMap) Delete(key string) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *UnAckedEventsMap) Store(key string, value models.Event) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}
