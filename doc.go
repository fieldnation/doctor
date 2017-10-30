package doctor

import (
	"errors"
	"fmt"
)

// HealthCheck performs a checkup and returns a bill of health report.
type HealthCheck func(b BillOfHealth) BillOfHealth

// Doctor represents a worker who will perform different
// types of health checks periodically.
type Doctor struct {
	cal *calendar
	closeNotify chan struct{}
}

// New returns a new doctor.
func New() *Doctor { return &Doctor{cal: newCalendar()} }

// Schedule a doctor appointment with a variety of options.
func (d *Doctor) Schedule(a Appointment, opts ...Option) error {

	if a.Name == "" {
		return errors.New("name not set for an appointment")
	}

	if a.HealthCheck == nil {
		return errors.New("appointment health check not assigned")
	}

	// ensure no duplicate appointment names exist
	if _, ok := d.cal.get(a.Name); ok {
		return fmt.Errorf("unable to schedule health check: %q already exists", a.Name)
	}

	// create a new appointment, and set its options
	appt := newAppt(a.Name, a.HealthCheck)
	for _, opt := range opts {
		opt(&appt.opts) // for now we don't check option errs
	}

	// append the appointment to the doctors calendar
	d.cal.set(appt)

	return nil
}

// Examine starts the series of health checks that were registered.
func (d *Doctor) Examine() (<-chan BillOfHealth, <-chan CloseNotify struct{}) {

	// range over each appointment and begin the exam
	c := d.cal.begin()

	// when the waitgroup finishes, close the channel
	go func() {
		d.cal.wait()
		close(c)
	}()

	// return the BillOfHealth recieving channel
	return c
}

// BillsOfHealth returns a list of bills of health.
func (d *Doctor) BillsOfHealth() []BillOfHealth {
	bills := []BillOfHealth{}
	for _, a := range d.cal.exams {
		bills = append(bills, a.get())
	}
	return bills
}

// Close sends a kill signal to all the long running healthchecks.
func (d *Doctor) Close() {
	d.cal.close()
}
