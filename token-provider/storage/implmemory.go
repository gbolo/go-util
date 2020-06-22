package storage

import "sync"

type MemoryStore struct {
	sync.RWMutex
	services map[string]*Service
}

func NewMemoryStore() Store {
	return &MemoryStore{
		services: make(map[string]*Service, 0),
	}
}

func (m *MemoryStore) AddService(description string) *Service {
	service := NewService(description)
	if service != nil {
		m.Lock()
		defer m.Unlock()
		if _, found := m.services[service.ID]; !found {
			m.services[service.ID] = service
			return service
		}
	}
	return nil
}

func (m *MemoryStore) AddServiceWithID(id, description string) *Service {
	service := NewServiceWithID(id, description)
	if service != nil {
		m.Lock()
		defer m.Unlock()
		if _, found := m.services[service.ID]; !found {
			m.services[service.ID] = service
			return service
		}
	}
	return nil
}

func (m *MemoryStore) RemoveService(id string) {
	m.Lock()
	defer m.Unlock()
	if _, found := m.services[id]; found {
		delete(m.services, id)
	}
}

func (m *MemoryStore) UpdateService(service *Service) {
	if service == nil || service.ID == "" {
		return
	}
	m.Lock()
	defer m.Unlock()
	if _, found := m.services[service.ID]; found {
		m.services[service.ID] = service
	}
}

func (m *MemoryStore) ListServices() map[string]*Service {
	m.RLock()
	defer m.RUnlock()
	return m.services
}

func (m *MemoryStore) GetServiceByID(id string) *Service {
	m.RLock()
	defer m.RUnlock()
	if _, found := m.services[id]; found {
		return m.services[id]
	}
	return nil
}

func (m *MemoryStore) ValidateApiKeyForServiceID(id, rawKey string) bool {
	service := m.GetServiceByID(id)
	if service != nil {
		return service.ValidateApiKey(rawKey)
	}
	return false
}
