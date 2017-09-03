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

	// WaitGroup manages HealthCheck results
	wg      sync.WaitGroup
	results chan BillOfHealth
}

// New returns a new doctor.
func New() *Doctor {
	return &Doctor{results: make(chan BillOfHealth)}
}

// Schedule a health check with some options, bascially a doctor appointment.
func (d *Doctor) Schedule(name string, h HealthCheck, opts ...Options) error {

	// ensure no duplicate health checks names exist
	for _, a := range d.appts {
		if a.result.name == name {
			return fmt.Errorf("unable to schedule health check: %q already exists", name)
		}
	}

	// create a new appointment
	a := &appointment{
		healthCheck: h,
		result: BillOfHealth{
			name:        name,
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

	if d.examining {
		if err := d.examine(a); err != nil {
			return err
		}
	}

	return nil
}

// Examine starts the series of health checks that were registered.
func (d *Doctor) Examine() (<-chan BillOfHealth, error) {

	d.examining = true

	// range over each appointment
	for _, appt := range d.appts {
		if err := d.examine(appt); err != nil {
			return nil, err
		}
	}

	go func() {
		d.wg.Wait()
		close(d.results)
	}()

	return d.results, nil
}

func (d *Doctor) examine(appt *appointment) error {

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
					d.results <- boh
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

	return nil
}

// Results returns a list of bills of health.
func (d *Doctor) Results() []BillOfHealth {
	boh := []BillOfHealth{}
	for _, a := range d.appts {
		boh = append(boh, a.get())
	}
	return boh
}
