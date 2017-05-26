package backend

import (
	"encoding/json"
	"sync"
	"sync/atomic"
)

// Backend ...
type Backend struct {
	Address   string `json:"address"`
	ConnCount int32  `json:"conn_count"`
}

// DcreCount ...
func (b *Backend) DcreCount() {
	atomic.AddInt32(&b.ConnCount, -1)
}

func (b *Backend) getCount() int32 {
	return atomic.LoadInt32(&b.ConnCount)
}

// Manager ...
type Manager struct {
	backends []Backend
	mutex    sync.RWMutex
}

// NewManager ...
func NewManager(addresses []string) *Manager {
	m := Manager{}
	for _, item := range addresses {
		m.backends = append(m.backends, Backend{Address: item})
	}
	return &m
}

// Get ...
func (m *Manager) Get() *Backend {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if len(m.backends) == 0 {
		return nil
	}
	result := &m.backends[0]
	for i := 1; i < len(m.backends); i++ {
		if m.backends[i].getCount() < result.getCount() {
			result = &m.backends[i]
		}
	}
	atomic.AddInt32(&result.ConnCount, 1)
	return result
}

// Add ...
func (m *Manager) Add(address string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.backends = append(m.backends, Backend{Address: address})
}

// Describe ...
func (m *Manager) Describe() []byte {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	b, err := json.Marshal(m.backends)
	if err != nil {
		panic(err)
	}
	return b
}
