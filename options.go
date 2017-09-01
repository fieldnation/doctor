package doctor

import "time"

// Options takes an option and returns an error.
type Options func(*options) error

type options struct {
	ttl      time.Duration
	interval time.Duration
	verbose  bool
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

// Verbose sets the verbose option.
func Verbose() Options {
	return func(o *options) error {
		o.verbose = true
		return nil
	}
}
