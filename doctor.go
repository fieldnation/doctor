package doctor

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Doctor represents a worker who will perform different
// types of health checks periodically.
type Doctor struct {
	appts     []*appointment
	examining bool
}

// Options takes an option and returns an error.
type Options func(*options) error

type options struct {
	ttl      time.Duration
	interval time.Duration
}

type appointment struct {
	healthCheck HealthCheck
	opts        options

	// mu protects the bill of health
	mu     sync.RWMutex
	result BillOfHealth
}

func (a *appointment) set(r BillOfHealth) {
	a.mu.Lock()
	a.result = r
	a.mu.Unlock()
}

func (a *appointment) get() BillOfHealth {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.result
}

// BillOfHealth describes the results of a doctor appointment.
type BillOfHealth struct {
	Name        string `json:"name"`
	Body        []byte `json:"body"`
	ContentType string `json:"content_type"`
}

// New returns a new doctor.
func New() *Doctor {
	return &Doctor{}
}

// TTL sets the Time to Live option value.
func TTL(ttl time.Duration) Options {
	return func(o *options) error {
		o.ttl = ttl
		return nil
	}
}

// Regularity sets the duration of how often the health check is executed.
func Regularity(interval time.Duration) Options {
	return func(o *options) error {
		o.interval = interval
		return nil
	}
}

// HealthCheck performs a checkup and returns a report.
type HealthCheck func() (body []byte, contentType string, err error)

// Schedule a health check with some options, bascially a doctor appointment.
func (d *Doctor) Schedule(h HealthCheck, name string, opts ...Options) error {

	// check if an examination is already underway,
	// if so do not allow further scheduling, dynamic
	// scheduling is not yet supported
	if d.examining {
		return errors.New("you can only schedule health checks before the docker begins examinations")
	}

	// create a new appointment
	a := &appointment{
		healthCheck: h,
		result: BillOfHealth{
			Name:        name,
			Body:        []byte(fmt.Sprintf("{\"%s\": \"no health check results\"", name)),
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
func (d *Doctor) Examine() (<-chan time.Time, error) {

	// officially start examinations, if we ever support
	// dynamic/concurrent scheduling and examining this
	// would be removed
	d.examining = true

	var wg sync.WaitGroup
	out := make(chan time.Time)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan time.Time) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	// range over each appointment
	for _, a := range d.appts {

		quit := make(chan struct{})

		// execute the appointment in a seperate goroutine
		go func(appt *appointment) {

			// set a new ticker based on the scheduled regularity
			ticker := time.NewTicker(appt.opts.interval)
			wg.Add(1)
			go output(ticker.C)

			for {
				select {
				case <-ticker.C:
					go func(appt *appointment) {
						body, contentType, err := appt.healthCheck()
						if err != nil {
							close(quit)
							fmt.Printf("log: error: %s\n", err)
							return
						}
						appt.set(BillOfHealth{appt.result.Name, body, contentType})
					}(appt)
				case <-quit:
					ticker.Stop()
					fmt.Println("Stopped the ticker!")
					return
				}
			}
		}(a)

		// if there is a TTL set, close the appointment at that time
		if a.opts.ttl > 0 {

			// it is acceptable to range over each appointment sequencially
			// for setup, but TTL requires a goroutine to keep things moving
			go func(appt *appointment) {
				time.Sleep(appt.opts.ttl)
				close(quit)
			}(a)
		}

	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out, nil
}

// Results returns a list of bills of health.
func (d *Doctor) Results() []BillOfHealth {
	boh := []BillOfHealth{}
	for _, a := range d.appts {
		boh = append(boh, a.get())
	}
	return boh
}
