package doctor

import (
	"sync"
	"time"
)

// Appointment describes a particular doctor appointment.
type Appointment struct {
	Name        string      `json:"name, omitempty"`
	HealthCheck HealthCheck `json:"healthcheck, omitempty"`
}

type appointment struct {
	name string
	hc   HealthCheck
	opts options

	// mu protects the bill of health
	mu     sync.RWMutex
	status status
	boh    BillOfHealth
}

type status struct {
	quit   chan struct{}
	closed bool
}

func newAppt(name string, hc HealthCheck) *appointment {
	return &appointment{
		name:   name,
		hc:     hc,
		status: status{},
		boh: BillOfHealth{
			name:        name,
			closeNotify: make(chan struct{}),
			Body:        []byte("{\"report\": \"no health check results\"}"),
			ContentType: "application/json",
		}}
}

func (a *appointment) get() BillOfHealth {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.boh
}

func (a *appointment) close() {
	close(a.boh.closeNotify)
}

// run executes a healthcheck scheduled by an appointment,
// run takes an BillOfHealth channel to send the result
// to and an optional callback as a convience
func (a *appointment) run() BillOfHealth {

	// since we are under mutex protection
	// we can directly reference the boh
	a.mu.Lock()
	defer a.mu.Unlock()

	// update the start time
	a.boh.start = time.Now()

	// pass the bill of health copy to the health check,
	// execute the health check, and overwrite the
	// bill of health copy with the new bill of health
	// values returned by the health check
	a.boh = a.hc(a.boh)

	// update the end time
	a.boh.end = time.Now()

	return a.boh
}
