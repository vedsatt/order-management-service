package patterns_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/patterns"
)

func TestRetry(t *testing.T) {
	maxRetries := 5
	baseDelay := 100 * time.Millisecond

	t.Run("success operation", func(t *testing.T) {
		successOperation := func() error {
			return nil
		}

		err := patterns.Retry(successOperation, maxRetries, baseDelay)
		if err != nil {
			t.Fatalf("expected nil err, but got: %v", err)
		}
	})

	t.Run("failed operation", func(t *testing.T) {
		var expectedError = fmt.Errorf("error with operation")
		failedOperation := func() error {
			return expectedError
		}

		err := patterns.Retry(failedOperation, maxRetries, baseDelay)
		if err != expectedError {
			t.Fatalf("expected error: %v, but got: %v", expectedError, err)
		}
	})

	t.Run("success after 3 retries operation", func(t *testing.T) {
		callCount := 0
		lateSuccess := func() error {
			if callCount == 3 {
				return nil
			}
			callCount++
			return fmt.Errorf("error")
		}

		start := time.Now()
		err := patterns.Retry(lateSuccess, maxRetries, baseDelay)
		if err != nil {
			t.Fatalf("expected nil err, but got: %v", err)
		}

		if callCount != 3 {
			t.Fatalf("expected 3 attempts, got %d", callCount)
		}

		elapsed := time.Since(start)
		expectedTime := 300 * time.Millisecond
		if elapsed < expectedTime {
			t.Fatalf("retry should take at least %v, but took %v", expectedTime, elapsed)
		}
	})
}

func TestTimeout(t *testing.T) {
	timeout := time.Millisecond * 500

	t.Run("success operation", func(t *testing.T) {
		successOperation := func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		}

		err := patterns.Timeout(successOperation, timeout)
		if err != nil {
			t.Fatalf("expected nil error, but got: %v", err)
		}
	})

	t.Run("success operation", func(t *testing.T) {
		expectedError := fmt.Errorf("error")
		errOperation := func() error {
			return expectedError
		}

		err := patterns.Timeout(errOperation, timeout)
		if err != expectedError {
			t.Fatalf("expected error: %v, but got: %v", expectedError, err)
		}
	})

	t.Run("success operation", func(t *testing.T) {
		timeoutExpiredOperation := func() error {
			time.Sleep(timeout * 2)
			return nil
		}

		expectedError := fmt.Errorf("waiting time exceeded")
		err := patterns.Timeout(timeoutExpiredOperation, timeout)
		if err == nil {
			t.Fatalf("expected err: %v, but got nil", expectedError)
		}
	})
}

func TestDLQ(t *testing.T) {
	t.Run("add messages to DLQ on failure", func(t *testing.T) {
		dlq := patterns.NewDeadLetterQueue()
		messages := []string{"msg1", "msg2", "msg3"}

		patterns.ProcessWithDLQ(messages, func(msg string) error {
			if msg == "msg2" {
				return fmt.Errorf("processing failed")
			}
			return nil
		}, dlq)

		if dlq.Size() != 1 {
			t.Errorf("expected 1 message in DLQ, got %d", dlq.Size())
		}

		if dlq.FailedAmount() != 1 {
			t.Errorf("expected 1 failed message, got %d", dlq.FailedAmount())
		}

		dlqMessages := dlq.GetMessages()
		if len(dlqMessages) != 1 || dlqMessages[0] != "msg2" {
			t.Errorf("expected ['msg2'] in DLQ, got %v", dlqMessages)
		}
	})

	t.Run("no messages in DLQ on success", func(t *testing.T) {
		dlq := patterns.NewDeadLetterQueue()
		messages := []string{"msg1", "msg2", "msg3"}

		patterns.ProcessWithDLQ(messages, func(msg string) error {
			return nil
		}, dlq)

		if dlq.Size() != 0 {
			t.Errorf("expected 0 messages in DLQ, got %d", dlq.Size())
		}
	})

	t.Run("clear DLQ", func(t *testing.T) {
		dlq := patterns.NewDeadLetterQueue()
		messages := []string{"msg1", "msg2"}

		patterns.ProcessWithDLQ(messages, func(msg string) error {
			return fmt.Errorf("always fail")
		}, dlq)

		if dlq.Size() != 2 {
			t.Errorf("expected 2 messages in DLQ before clear")
		}

		dlq.Clear()

		if dlq.Size() != 0 {
			t.Errorf("expected 0 messages in DLQ after clear, got %d", dlq.Size())
		}

		if dlq.FailedAmount() != 2 {
			t.Errorf("expected failed counter reset after clear, got %d", dlq.FailedAmount())
		}
	})

	t.Run("thread safety", func(t *testing.T) {
		dlq := patterns.NewDeadLetterQueue()
		var wg sync.WaitGroup

		for i := range 50 {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				dlq.Add(fmt.Sprintf("msg%d", id))
			}(i)
		}

		wg.Wait()

		if dlq.Size() != 50 {
			t.Errorf("expected 100 messages in DLQ, got %d", dlq.Size())
		}

		if dlq.FailedAmount() != 50 {
			t.Errorf("expected 100 failed messages, got %d", dlq.FailedAmount())
		}
	})
}
