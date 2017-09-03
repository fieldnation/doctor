package doctor

import "time"

// BillOfHealth describes the results of a doctor appointment.
type BillOfHealth struct {
	err         error
	start       time.Time
	end         time.Time
	name        string
	healthy     bool
	Body        []byte `json:"body"`
	ContentType string `json:"content_type"`
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

// Err returns a BillOfHealth Err
func (b BillOfHealth) Err() error {
	return b.err
}

// Start returns the start of BillOfHealth Timestamp
func (b BillOfHealth) Start() time.Time {
	return b.start
}

// End returns a BillOfHealth Timestamp
func (b BillOfHealth) End() time.Time {
	return b.end
}
