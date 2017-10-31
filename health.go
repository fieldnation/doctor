package doctor

import "time"

// Health describes the results of a doctor appointment.
type Health struct {
	name        string
	status      bool
	start       time.Time
	end         time.Time
	Body        []byte `json:"body, omitempty"`
	ContentType string `json:"content_type, omitempty"`
	closeNotify chan struct{}
}

// Name returns the name of the Health.
func (h Health) Name() string {
	return h.name
}

// Healthy reports the current health status.
func (h Health) Healthy() bool {
	return h.status
}

// SetHealthy sets the health status to true.
func (h Health) SetHealthy() {
	h.status = true
}

// SetUnhealthy sets the health status to true.
func (h Health) SetUnhealthy() {
	h.status = false
}

// Start returns the start of BillOfHealth Timestamp
func (h Health) Start() time.Time {
	return h.start
}

// End returns a BillOfHealth Timestamp
func (h Health) End() time.Time {
	return h.end
}

// CloseNotify returns an channel that recieves empty structs.
// When the appointment closes that channel will recieve an
// empty struct.
func (h Health) CloseNotify() <-chan struct{} {
	return h.closeNotify
}
