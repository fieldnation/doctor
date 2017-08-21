package doctor

// Doctor represents a worker who will perform different
// types of health checks periodically.
type Doctor struct {
	opts  options
	appts []appt
}

// Options takes an option and returns an error.
type Options func(*options) error

type options struct {
	ttl      int
	interval int
}

type appt struct {
	healthCheck  HealthCheck
	opts         options
	billOfHealth []byte
}

// New returns a new doctor.
func New() *Doctor {
	return &Doctor{}
}

// TTL sets the Time to Live option value.
func TTL(ttl int) Options {
	return func(o *options) error {
		o.ttl = ttl
		return nil
	}
}

// Interval sets the Interval option value.
func Interval(interval int) Options {
	return func(o *options) error {
		o.interval = interval
		return nil
	}
}

// HealthCheck performs a checkup and returns a report.
type HealthCheck func() (body []byte, contentType string, err error)

// Schedule a health check with some options, bascially a doctor appointment.
func (d *Doctor) Schedule(h HealthCheck, opts ...Options) {

	// create a new appointment
	a := appt{
		healthCheck:  h,
		billOfHealth: []byte("{\"report\": \"no health check results\""),
	}

	// set the request options on that appointment
	for _, o := range opts { // for now we don't check err
		o(&a.opts)
	}

	// append the appointment to the doctors list
	d.appts = append(d.appts, a)
}

// Examine starts the series of health checks that were registered.
func (d *Doctor) Examine() (<-chan bool, error) {
	return nil, nil
}
