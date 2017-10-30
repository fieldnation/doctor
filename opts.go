package doctor

import "time"

// Option takes an option and returns an error.
type Option func(*options) error

type options struct {
	ttl      time.Duration
	interval time.Duration
	sleep    time.Duration
	attempts int
	verbose  bool
}

// Verbose sets the verbose option.
func Verbose() Option {
	return func(o *options) error {
		o.verbose = true
		return nil
	}
}

// Regularity sets the duration of how often the health check is executed.
func Regularity(interval time.Duration) Option {
	return func(o *options) error {
		o.interval = interval
		return nil
	}
}

// TTL sets the Time to Live option value.
func TTL(ttl time.Duration) Option {
	return func(o *options) error {
		o.ttl = ttl
		return nil
	}
}

// MaxAttempts sets the maximum number of repeated failures
// before closing the health check.
func MaxAttempts(attempts int) Option {
	return func(o *options) error {
		o.attempts = attempts
		return nil
	}
}

// Sleep waits a period of time before starting
// health check execution.
func Sleep(sleep time.Duration) Option {
	return func(o *options) error {
		o.sleep = sleep
		return nil
	}
}
