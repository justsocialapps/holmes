// Package justlib contains handy utility functions to make your day-to-day Go
// programming more pleasant.
package justlib

import "time"

// Try calls the given function »f« at most »times« times as long as it returns
// an error, sleeping the given »backoffTime« in between subsequent calls.  If
// f returns an error on the last try, that error is returned by Try.
// If times is 0, then f is called indefinitely as long as it returns an error.
func Try(times uint16, backoffTime time.Duration, f func() error) error {
	var err error
	var i uint16
	for i = 1; ; i++ {
		err = f()
		if err == nil {
			return nil
		}
		if times > 0 && i >= times {
			return err
		}
		time.Sleep(backoffTime)
	}
}
