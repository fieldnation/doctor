package doctor

import "sync"

type appointment struct {
	healthCheck HealthCheck
	opts        options

	// mu protects the bill of health
	mu     sync.RWMutex
	result BillOfHealth
}

func (a *appointment) set(r BillOfHealth) {
	a.mu.Lock()
	a.result = r
	a.mu.Unlock()
}

func (a *appointment) get() BillOfHealth {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.result
}
