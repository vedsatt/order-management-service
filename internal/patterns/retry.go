package patterns

import (
	"time"
)

func Retry(operation func() error, maxRetries int, baseDelay time.Duration) error {
	var err error

	for attempt := range maxRetries {
		if err = operation(); err == nil {
			return nil
		}

		if attempt < maxRetries-1 {
			delay := baseDelay * (1 << attempt)
			time.Sleep(delay)
		}
	}

	return err
}
