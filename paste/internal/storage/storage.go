package storage

import "fmt"

type Storage interface {
	Save(url, content string) error
	Get(url string) (string, error)
	Exists(url string) bool
	Delete(url string) error
	List() []string
}

type MemoryStorage struct {
	pastes map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		pastes: make(map[string]string),
	}
}

func (m *MemoryStorage) Save(url, content string) error {
	if _, exists := m.pastes[url]; exists {
		return fmt.Errorf("URL already exists")
	}
	m.pastes[url] = content
	return nil
}

func (m *MemoryStorage) Get(url string) (string, error) {
	val, exists := m.pastes[url]
	if !exists {
		return "", fmt.Errorf("URL not found")
	}
	return val, nil
}

func (m *MemoryStorage) Exists(url string) bool {
	_, exists := m.pastes[url]
	return exists
}

func (m *MemoryStorage) Delete(url string) error {
	if _, exists := m.pastes[url]; exists {
		delete(m.pastes, url)
		return nil
	}
	return fmt.Errorf("URL not found")
}

func (m *MemoryStorage) List() []string {
	var urls []string
	for key := range m.pastes {
		urls = append(urls, key)
	}
	return urls
}
