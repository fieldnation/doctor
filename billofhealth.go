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

// Healthy() returns healthy
func (b BillOfHealth) Healthy() error {
	b.healthy = true
	return nil
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
