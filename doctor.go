package doctor

import (
	"fmt"
	"sync"
	"time"
)

// HealthCheck performs a checkup and returns a bill of health report.
type HealthCheck func(b BillOfHealth) BillOfHealth

// Doctor represents a worker who will perform different
// types of health checks periodically.
type Doctor struct {
	calendar []*appointment

	r *reaper

	// wg waits on HealthChecks complete
	// scheduled executions or the BillOfHealth
	// channel to finish draining
	wg sync.WaitGroup
	c  chan BillOfHealth
}

// New returns a new doctor.
func New() *Doctor { return &Doctor{r: newReaper()} }

// Schedule a doctor appointment with a variety of options.
func (d *Doctor) Schedule(appt Appointment, opts ...Option) error {

	// ensure no duplicate health check names exist
	for _, a := range d.calendar {
		if a.boh.name == appt.Name {
			return fmt.Errorf("unable to schedule health check: %q already exists", appt.Name)
		}
	}

	// create a new appointment, and set its options
	a := newAppt(appt.Name, appt.HealthCheck)
	for _, o := range opts {
		o(&a.opts) // for now we don't check option errs
	}

	// append the appointment to the doctors list
	d.calendar = append(d.calendar, a)

	return nil
}

// Examine starts the series of health checks that were registered.
func (d *Doctor) Examine() <-chan BillOfHealth {

	// create a BillOfHealth channel with a buffer equal to
	// the number of appointments scheduled on the calendar
	d.c = make(chan BillOfHealth, len(d.calendar))

	// range over each appointment and begin the exam
	for _, appt := range d.calendar {
		d.examine(appt)
	}

	// when the waitgroup finishes, close the channel
	go func() {
		d.wg.Wait()
		close(d.c)
		d.c = nil
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

	go func() {
		if interval < 1 {
			d.run(appt, func() { d.wg.Done() })
			return
		}
		tick := time.NewTicker(interval)
		done := make(chan struct{})
		if err := d.r.Set(appt.name, done); err != nil {
			return // do not execute the ticker
		}
		for {
			select {
			case <-tick.C:
				go d.run(appt)
			case <-done:
				tick.Stop()
				d.wg.Done()
				appt.close()
				return
			}
		}
	}()

	// if a TTL is set, close the channel at that time
	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			done, ok := d.r.Get(appt.name)
			if ok {
				close(done)
				d.r.Delete(appt.name)
			}
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

// Close sends a kill signal to all the long running healthchecks.
func (d *Doctor) Close() {
	d.r.Close()
}

// run executes a healthcheck scheduled by an appointment,
// run takes an appointment and an optional callback,
// if you don't need a callback, simply pass a nil value
// as the second parameter
func (d *Doctor) run(appt *appointment, callbacks ...func()) {

	// send the bill of health result down the output channel
	if d.c == nil {
		fmt.Printf("d.c is nil: %q\n", appt.name)
		return
	}
	d.c <- appt.run()

	// execute callbacks
	for _, f := range callbacks {
		f()
	}
}
