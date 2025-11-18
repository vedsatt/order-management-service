package patterns

import "sync"

type DeadLetterQueue struct {
	queue  []string
	mu     sync.RWMutex
	failed int
}

func NewDeadLetterQueue() *DeadLetterQueue {
	return &DeadLetterQueue{
		queue:  make([]string, 0),
		mu:     sync.RWMutex{},
		failed: 0,
	}
}

func (l *DeadLetterQueue) GetMessages() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.queue
}

func (l *DeadLetterQueue) Add(msg ...string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.failed++
	l.queue = append(l.queue, msg...)
}

func (l *DeadLetterQueue) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = make([]string, 0)
}

func (l *DeadLetterQueue) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.queue)
}

func (l *DeadLetterQueue) FailedAmount() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.failed
}

func ProcessWithDLQ(messages []string, operation func(msg string) error, dlq *DeadLetterQueue) {
	for _, message := range messages {
		err := operation(message)
		if err != nil {
			dlq.Add(message)
		}
	}
}
