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
	appts     []*appointment
	examining bool

	// WaitGroup manages HealthCheck BillOfHealth channel
	wg sync.WaitGroup
	c  chan BillOfHealth
}

// New returns a new doctor.
func New() *Doctor {
	return &Doctor{c: make(chan BillOfHealth)}
}

// Schedule a health check with some options, bascially a doctor appointment.
func (d *Doctor) Schedule(appt Appointment, opts ...Options) error {

	// ensure no duplicate health checks names exist
	for _, a := range d.appts {
		if a.boh.name == appt.Name {
			return fmt.Errorf("unable to schedule health check: %q already exists", appt.Name)
		}
	}

	// create a new appointment
	a := newAppt(appt.Name, appt.HealthCheck)

	// set the request options on that appointment
	for _, o := range opts { // for now we don't check err
		o(&a.opts)
	}

	// append the appointment to the doctors list
	d.appts = append(d.appts, a)

	// if Examine() has been called, then have the
	// appointment begin the examination
	if d.examining {
		d.examine(a)
	}

	return nil
}

// Examine starts the series of health checks that were registered.
func (d *Doctor) Examine() <-chan BillOfHealth {

	d.examining = true

	// range over each appointment and begin the exam
	for _, appt := range d.appts {
		d.examine(appt)
	}

	// wait for the waitgroup to finish, then close the channel
	go func() {
		d.wg.Wait()
		close(d.c)
	}()

	return d.c
}

func (d *Doctor) examine(appt *appointment) {

	// add to the WaitGroup before starting the
	// goroutine to avoid wg.Wait() returning
	// before wg.Add(1) can be executed
	d.wg.Add(1)
	done := make(chan struct{}) // close channel for the ticker

	// execute the appointment in a seperate goroutine
	go func(app *appointment, quit chan struct{}) {

		// if the interval is less than one, simply
		// execute the health check once
		if app.opts.interval < 1 {
			go d.run(appt, func() {
				d.wg.Done()
			})
			return
		}

		// create a ticker at the requested interval rate
		// and execute the health check at every tick
		ticker := time.NewTicker(app.opts.interval)
		for {
			select {
			case <-ticker.C:
				go func(a *appointment) {
					go d.run(appt, nil)
				}(app)
			case <-quit: // quit signal to stop the ticker
				ticker.Stop()
				d.wg.Done()
				return
			}
		}
	}(appt, done)

	// if there is a TTL set, close the appointment at that time
	if appt.opts.ttl > 0 {
		go func(app *appointment, quit chan struct{}) {
			<-time.After(app.opts.ttl)
			close(quit)
		}(appt, done)
	}
}

// BillsOfHealth returns a list of bills of health.
func (d *Doctor) BillsOfHealth() []BillOfHealth {
	bills := []BillOfHealth{}
	for _, a := range d.appts {
		bills = append(bills, a.get())
	}
	return bills
}

func (d *Doctor) run(appt *appointment, callback func()) {
	boh := appt.get()
	boh.start = time.Now()
	boh = appt.healthCheck(boh)
	boh.end = time.Now()
	appt.set(boh)
	d.c <- boh
	if callback != nil {
		callback()
	}
}
