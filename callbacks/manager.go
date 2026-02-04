package callbacks

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
)

var ErrCallbackNotFound = errors.New("callback not found")

// Repository defines callback storage.
type Repository interface {
	Create(syncID uuid.UUID) error
	Set(callback models.MethodCallback) error
	Wait(ctx context.Context, syncID uuid.UUID) (models.MethodCallback, error)
	Pop(syncID uuid.UUID) (models.MethodCallback, bool)
	Stop() error
}

// MemoryRepo stores callbacks in-memory.
type MemoryRepo struct {
	mu        sync.Mutex
	waiting   map[uuid.UUID]chan models.MethodCallback
	callbacks map[uuid.UUID]models.MethodCallback
	stopped   bool
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{waiting: make(map[uuid.UUID]chan models.MethodCallback), callbacks: make(map[uuid.UUID]models.MethodCallback)}
}

func (r *MemoryRepo) Create(syncID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.stopped {
		return context.Canceled
	}
	if _, ok := r.waiting[syncID]; !ok {
		r.waiting[syncID] = make(chan models.MethodCallback, 1)
	}
	return nil
}

func (r *MemoryRepo) Set(callback models.MethodCallback) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.stopped {
		return context.Canceled
	}
	ch, ok := r.waiting[callback.SyncID]
	if ok {
		ch <- callback
		return nil
	}
	r.callbacks[callback.SyncID] = callback
	return nil
}

func (r *MemoryRepo) Wait(ctx context.Context, syncID uuid.UUID) (models.MethodCallback, error) {
	r.mu.Lock()
	if cb, ok := r.callbacks[syncID]; ok {
		delete(r.callbacks, syncID)
		r.mu.Unlock()
		return cb, nil
	}
	ch, ok := r.waiting[syncID]
	if !ok {
		ch = make(chan models.MethodCallback, 1)
		r.waiting[syncID] = ch
	}
	r.mu.Unlock()

	select {
	case cb := <-ch:
		return cb, nil
	case <-ctx.Done():
		return models.MethodCallback{}, ctx.Err()
	}
}

func (r *MemoryRepo) Pop(syncID uuid.UUID) (models.MethodCallback, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cb, ok := r.callbacks[syncID]
	if ok {
		delete(r.callbacks, syncID)
	}
	return cb, ok
}

func (r *MemoryRepo) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stopped = true
	for _, ch := range r.waiting {
		close(ch)
	}
	return nil
}

// Manager coordinates callback waiting and timeouts.
type Manager struct {
	repo Repository
}

func NewManager(repo Repository) *Manager {
	return &Manager{repo: repo}
}

func (m *Manager) Create(syncID uuid.UUID) error {
	return m.repo.Create(syncID)
}

func (m *Manager) Set(callback models.MethodCallback) error {
	return m.repo.Set(callback)
}

func (m *Manager) Wait(syncID uuid.UUID, timeout time.Duration) (models.MethodCallback, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return m.repo.Wait(ctx, syncID)
}

func (m *Manager) Stop() error {
	return m.repo.Stop()
}
