package doctor

import "sync"

// Appointment describes a particular doctor appointment.
type Appointment struct {
	Name        string
	HealthCheck HealthCheck
}

type appointment struct {
	healthCheck HealthCheck
	opts        options

	// mu protects the bill of health
	mu  sync.RWMutex
	boh BillOfHealth
}

func (a *appointment) set(boh BillOfHealth) {
	a.mu.Lock()
	a.boh = boh
	a.mu.Unlock()
}

func (a *appointment) get() BillOfHealth {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.boh
}
