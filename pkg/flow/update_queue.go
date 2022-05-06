package flow

import (
	"sync"
)

// updateQueue is an unordered queue of updated components to be processed.
type updateQueue struct {
	mut     sync.Mutex
	updated map[*userComponent]struct{}

	updateCh chan struct{}
}

func newUpdateQueue() *updateQueue {
	return &updateQueue{
		updateCh: make(chan struct{}, 1),
		updated:  make(map[*userComponent]struct{}),
	}
}

// Enqueue enqueues a new userComponent.
func (uq *updateQueue) Enqueue(uc *userComponent) {
	uq.mut.Lock()
	uq.updated[uc] = struct{}{}
	uq.mut.Unlock()

	select {
	case uq.updateCh <- struct{}{}:
	default:
	}
}

// UpdateCh returns a channel which will return a value when Dequeue can be
// called.
func (uq *updateQueue) UpdateCh() <-chan struct{} { return uq.updateCh }

// TryDequeue tries to dequeue a userComponent. Returns nil if the queue is
// empty.
func (uq *updateQueue) TryDequeue() *userComponent {
	uq.mut.Lock()
	defer uq.mut.Unlock()

	for uc := range uq.updated {
		delete(uq.updated, uc)
		return uc
	}

	return nil
}