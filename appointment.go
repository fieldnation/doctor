package doctor

import "sync"

// Appointment describes a particular doctor appointment.
type Appointment struct {
	Name        string      `json:"name, omitempty"`
	HealthCheck HealthCheck `json:"healthcheck, omitempty"`
}

type appointment struct {
	hc   HealthCheck
	opts options

	// mu protects the bill of health
	mu  sync.RWMutex
	boh BillOfHealth
}

func newAppt(name string, hc HealthCheck) *appointment {
	return &appointment{hc: hc, boh: BillOfHealth{
		name:        name,
		Body:        []byte("{\"report\": \"no health check results\"}"),
		ContentType: "application/json",
	}}
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
