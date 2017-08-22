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
	appts     []appt
	examining bool
}

// Options takes an option and returns an error.
type Options func(*options) error

type options struct {
	ttl      time.Duration
	interval time.Duration
}

type appt struct {
	healthCheck HealthCheck
	opts        options

	mu           sync.RWMutex
	billOfHealth []byte
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
func (d *Doctor) Schedule(h HealthCheck, opts ...Options) error {

	// check if an examination is already underway,
	// if so do not allow further scheduling, dynamic
	// scheduling is not yet supported
	if d.examining {
		return errors.New("you can only schedule health checks before the docker begins examinations")
	}

	// create a new appointment
	a := appt{
		healthCheck:  h,
		billOfHealth: []byte("{\"report\": \"no health check results\""),
		opts:         options{interval: 5 * time.Second},
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
func (d *Doctor) Examine() (<-chan bool, error) {

	// officially start examinations, if we ever support
	// dynamic/concurrent scheduling and examining this
	// would be removed
	d.examining = true

	// range over each appointment
	for k := range d.appts {

		// set a new ticker based on the scheduled regularity
		ticker := time.NewTicker(d.appts[k].opts.interval)
		quit := make(chan struct{})

		// execute the appointment in a seperate goroutine
		go func() {
			for {
				select {
				case <-ticker.C:
					go func() {
						d.appts[k].healthCheck()
					}()
				case <-quit:
					ticker.Stop()
					fmt.Println("Stopped the ticker!")
					return
				}
			}
		}()

		// if there is a TTL set, close the appointment at that time
		if d.appts[k].opts.ttl > 0 {

			// it is acceptable to range over each appointment sequencially
			// for setup, but TTL requires a goroutine to keep things moving
			go func() {
				time.Sleep(d.appts[k].opts.ttl)
				close(quit)
			}()
		}

	}

	return nil, nil
}
