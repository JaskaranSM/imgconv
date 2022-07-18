package storage

import (
	"imgconv/errors"
	"sync"
)

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		images: make(map[string][]byte),
	}
}

type MemoryStorage struct {
	mut    sync.Mutex
	images map[string][]byte
}

func (m *MemoryStorage) GetImageBytesById(id string) ([]byte, error) {
	m.mut.Lock()
	defer m.mut.Unlock()
	for i, data := range m.images {
		if i == id {
			return data, nil
		}
	}
	return nil, errors.ErrorImageNotFound
}

func (m *MemoryStorage) SaveImageBytesId(id string, data []byte) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.images[id] = data
	return nil
}
