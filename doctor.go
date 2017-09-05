package doctor

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// HealthCheck performs a checkup and returns a bill of health report.
type HealthCheck func(b BillOfHealth) BillOfHealth

// Doctor represents a worker who will perform different
// types of health checks periodically.
type Doctor struct {
	calendar  []*appointment
	examining bool

	// wg waits on HealthChecks complete
	// scheduled executions or the BillOfHealth
	// channel to finish draining
	wg sync.WaitGroup
	c  chan BillOfHealth
}

// New returns a new doctor.
func New() *Doctor { return &Doctor{} }

// Schedule a doctor appointment with a variety of options.
func (d *Doctor) Schedule(appt Appointment, opts ...Option) error {

	// since we don't have a strategy for concurrent/parallel
	// scheduling during examination, lets remove that senario
	// all together
	if d.examining {
		return errors.New("you cannot schedule an appointment if examinations have already begun")
	}

	// ensure no duplicate health check names exist
	for _, a := range d.calendar {
		if a.boh.name == appt.Name {
			return fmt.Errorf("unable to schedule health check: %q already exists", appt.Name)
		}
	}

	// create a new appointment
	a := newAppt(appt.Name, appt.HealthCheck)

	// set the request options on that appointment
	for _, o := range opts {
		o(&a.opts) // for now we don't check option errs
	}

	// append the appointment to the doctors list
	d.calendar = append(d.calendar, a)

	// if Examine() has been called, then have the
	// appointment begin the examination
	if d.examining {
		d.examine(a)
	}

	return nil
}

// Examine starts the series of health checks that were registered.
func (d *Doctor) Examine() <-chan BillOfHealth {

	// set examination state
	d.examining = true

	// make a BillOfHealth channel with a buffer equal
	// to the number of appointments scheduled on the
	// doctors calendar
	d.c = make(chan BillOfHealth, len(d.calendar))

	// range over each appointment and begin the exam
	for _, appt := range d.calendar {
		d.examine(appt)
	}

	// when the waitgroup finishes, close the channel
	go func() {
		d.wg.Wait()
		close(d.c)
	}()

	// return the BillOfHealth recieving channel
	return d.c
}

func (d *Doctor) examine(appt *appointment) {

	// set state
	interval := appt.opts.interval
	ttl := appt.opts.ttl

	// add to the WaitGroup before starting the
	// goroutine to avoid wg.Wait() returning
	// before wg.Add(1) can be executed
	d.wg.Add(1)

	// if the interval is less than one, simply
	// execute the health check once in a seperate
	// goroutine
	if interval < 1 {
		go d.run(appt, func() { d.wg.Done() })
		return
	}

	// quit channel for the ticker
	done := make(chan struct{})
	go func() {
		// create a ticker at the requested interval rate
		// and execute the health check at every tick
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				go d.run(appt)
			case <-done:
				ticker.Stop()
				d.wg.Done()
				return
			}
		}
	}()

	// if a TTL is set, close the channel at that time
	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			close(done)
		}()
	}
}

// BillsOfHealth returns a list of bills of health.
func (d *Doctor) BillsOfHealth() []BillOfHealth {
	bills := []BillOfHealth{}
	for _, a := range d.calendar {
		bills = append(bills, a.get())
	}
	return bills
}

// run executes a healthcheck scheduled by an appointment,
// run takes an appointment and an optional callback,
// if you don't need a callback, simply pass a nil value
// as the second parameter
func (d *Doctor) run(appt *appointment, callbacks ...func()) {

	// get a copy, (not a pointer) of the latest
	// bill of health in a thread safe manner
	boh := appt.get()

	// update the start time
	boh.start = time.Now()

	// pass the bill of health copy to the health check,
	// execute the health check, and overwrite the
	// bill of health copy with the new bill of health
	// values returned by the health check
	boh = appt.hc(boh)

	// update the end time
	boh.end = time.Now()

	// update the appointment with the new
	// bill of health result, in a thread safe manner
	appt.set(boh)

	// send the bill of health result down the output channel
	d.c <- boh

	// execute callbacks
	for _, f := range callbacks {
		f()
	}
}
