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
	appts     []*appointment
	examining bool
}

// New returns a new doctor.
func New() *Doctor {
	return &Doctor{}
}

// Schedule a health check with some options, bascially a doctor appointment.
func (d *Doctor) Schedule(name string, h HealthCheck, opts ...Options) error {

	// check if an examination is already underway,
	// if so do not allow further scheduling, dynamic
	// scheduling is not yet supported
	if d.examining {
		panic(errors.New("you can only schedule health checks before the docker begins examinations"))
	}

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

	return nil
}

// Examine starts the series of health checks that were registered.
func (d *Doctor) Examine() (<-chan BillOfHealth, error) {

	// officially start examinations, if we ever support
	// dynamic/concurrent scheduling and examining this
	// would be removed
	d.examining = true

	// waitgoup manages HealthCheck results
	var wg sync.WaitGroup
	results := make(chan BillOfHealth)
	sync := make(chan struct{})

	// range over each appointment
	for _, apt := range d.appts {

		// create a close channel
		quit := make(chan struct{})

		// execute the appointment in a seperate goroutine
		go func(a *appointment) {
			wg.Add(1)
			ticker := time.NewTicker(a.opts.interval) // ticker set on regularity
			close(sync)
			for {
				select {
				// execute the healthcheck with every tick
				case t := <-ticker.C:
					fmt.Println("recieved tick")
					go func(a *appointment) {
						boh := a.healthCheck(BillOfHealth{start: t})
						a.set(boh)
						results <- boh
					}(apt)
				// listen for the quit signal to stop the ticker
				case <-quit:
					ticker.Stop()
					wg.Done()
					return
				}
			}
		}(apt)

		// if there is a TTL set, close the appointment at that time
		if apt.opts.ttl > 0 {

			// it is acceptable to range over each appointment sequencially
			// for setup, but TTL requires a goroutine to keep things moving
			go func(a *appointment) {
				time.Sleep(a.opts.ttl)
				close(quit)
			}(apt)
		}
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	<-sync
	go func() {
		wg.Wait()
		close(results)
	}()

	return results, nil
}

// Results returns a list of bills of health.
func (d *Doctor) Results() []BillOfHealth {
	boh := []BillOfHealth{}
	for _, a := range d.appts {
		boh = append(boh, a.get())
	}
	return boh
}
