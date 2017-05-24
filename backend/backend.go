package backend

import "sync/atomic"

// Backend ...
type Backend struct {
	Address   string
	connCount int32
}

// Manager ...
type Manager struct {
	backends []Backend
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
	if len(m.backends) == 0 {
		return nil
	}
	result := &m.backends[0]
	for i := 1; i < len(m.backends); i++ {
		if m.backends[i].connCount < result.connCount {
			result = &m.backends[i]
		}
	}
	atomic.AddInt32(&result.connCount, 1)
	return result
}

// Decr ...
func (b *Backend) Decr() {
	atomic.AddInt32(&b.connCount, -1)
}
