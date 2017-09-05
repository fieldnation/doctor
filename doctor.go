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
	a := &appointment{
		healthCheck: appt.HealthCheck,
		boh: BillOfHealth{
			name:        appt.Name,
			Body:        []byte("{\"report\": \"no health check results\"}"),
			ContentType: "application/json",
		},
		opts: options{interval: 5 * time.Second},
	}

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
	quit := make(chan struct{}) // close channel for the ticker

	// execute the appointment in a seperate goroutine
	go func(app *appointment, done chan struct{}) {
		ticker := time.NewTicker(app.opts.interval) // ticker set by regularity interval
		for {
			select {
			// execute the healthcheck with every tick
			case t := <-ticker.C:
				go func(a *appointment) {
					boh := a.get()
					boh.start = t
					boh = a.healthCheck(boh)
					a.set(boh)
					d.c <- boh
				}(app)
			// listen for the quit signal to stop the ticker
			case <-done:
				ticker.Stop()
				d.wg.Done()
				return
			}
		}
	}(appt, quit)

	// if there is a TTL set, close the appointment at that time
	if appt.opts.ttl > 0 {
		go func(app *appointment, done chan struct{}) {
			<-time.After(app.opts.ttl)
			close(done)
		}(appt, quit)
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
