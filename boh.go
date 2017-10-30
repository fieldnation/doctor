package doctor

import "time"

// BillOfHealth describes the results of a doctor appointment.
type BillOfHealth struct {
	name        string
	healthy     bool
	start       time.Time
	end         time.Time
	closeNotify chan struct{}
	Body        []byte `json:"body, omitempty"`
	ContentType string `json:"content_type, omitempty"`
}

// Name returns the name of the BillOfHealth.
func (b BillOfHealth) Name() string {
	return b.name
}

// Healthy sets the BillOfHealth to a healthy state.
func (b BillOfHealth) Healthy() bool {
	return b.healthy
}

// SetHealth sets the BillOfHealth healthy value.
func (b BillOfHealth) SetHealth(health bool) {
	b.healthy = health
}

// Start returns the start of BillOfHealth Timestamp
func (b BillOfHealth) Start() time.Time {
	return b.start
}

// End returns a BillOfHealth Timestamp
func (b BillOfHealth) End() time.Time {
	return b.end
}

// CloseNotify returns an channel that recieves empty structs.
// When the appointment closes that channel will recieve an
// empty struct.
func (b BillOfHealth) CloseNotify() <-chan struct{} {
	return b.closeNotify
}
