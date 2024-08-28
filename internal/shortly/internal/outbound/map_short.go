package outbound

import (
	"context"
	"sync"

	"github.com/shandysiswandi/gostarter/internal/shortly/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/codec"
)

type MapShort struct {
	json codec.Codec
	data map[string]domain.Short
	mu   sync.RWMutex
}

func NewMapShort(json codec.Codec) *MapShort {
	return &MapShort{json: json, data: make(map[string]domain.Short)}
}

func (ms *MapShort) Get(_ context.Context, key string) (*domain.Short, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	short, ok := ms.data[key]
	if !ok {
		return nil, nil //nolint:nilnil // no rows is not an error, just a nil result
	}

	return &short, nil
}

func (ms *MapShort) Set(_ context.Context, value domain.Short) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.data[value.Key] = value

	return nil
}
