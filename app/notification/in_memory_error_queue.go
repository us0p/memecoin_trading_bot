package notification

import (
	"time"
)

func newInMemoryErrorQueueKey(mint string, workflow Workflow) InMemoryErrorQueueKey {
	return InMemoryErrorQueueKey{
		mint,
		workflow,
	}
}

func newErrorNotification(err error, sev Severity) ErrorNotification {
	return ErrorNotification{
		Err:         err,
		Sent:        false,
		StartedAt:   time.Now(),
		ErrSeverity: sev,
	}
}
