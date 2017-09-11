package doctor

import (
	"fmt"
	"sync"
)

type reaper struct {
	mu      sync.RWMutex
	workers map[string]worker
}

type worker struct {
	quit   chan struct{}
	closed bool
}

func newReaper() *reaper {
	return &reaper{workers: make(map[string]worker)}
}

func (r *reaper) Set(name string, done chan struct{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.workers[name] = worker{quit: done}
	return nil
}

func (r *reaper) Delete(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.workers, name)
}

func (r *reaper) Get(name string) (chan struct{}, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	w, ok := r.workers[name]
	return w.quit, ok
}

func (r *reaper) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for name, worker := range r.workers {
		if worker.closed {
			continue
		}
		fmt.Printf("closing %q\n", name)
		close(worker.quit)
	}
}
