package storage

import "ai-pdb/internal/core"

type MemoryStore struct {
	data map[string]core.Document
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]core.Document),
	}
}

func (m *MemoryStore) Save(doc core.Document) {
	m.data[doc.ID] = doc
}

func (m *MemoryStore) List() []core.Document {
	result := make([]core.Document, 0, len(m.data))
	for _, v := range m.data {
		result = append(result, v)
	}
	return result
}
