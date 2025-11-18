package patterns

import (
	"context"
	"fmt"
	"time"
)

func Timeout(operation func() error, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errChan := make(chan error)
	go func() {
		errChan <- operation()
	}()

	select {
	case err := <-errChan:
		return err

	case <-ctx.Done():
		return fmt.Errorf("waiting time exceeded")
	}
}
