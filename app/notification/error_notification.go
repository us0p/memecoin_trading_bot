package notification

import (
	"fmt"
	"memecoin_trading_bot/app/constants"
	"time"
)

type ErrorReport struct {
	InMemoryErrorQueueKey
	ErrorNotification
}

func (n *Notifications) getRelevantErrors() []ErrorReport {
	reports := make([]ErrorReport, 0)

	for key, errs := range n.ErrQueue {
		for _, err := range errs {
			if !err.Sent && (err.ErrSeverity >= Core || isLongRunningError(err.StartedAt)) {
				err.Sent = true
				reports = append(reports, ErrorReport{
					InMemoryErrorQueueKey: key,
					ErrorNotification:     err,
				})
			}
		}
	}

	return reports
}

func (n *Notifications) errReport() []string {
	relevantErrors := n.getRelevantErrors()

	reports := make([]string, len(relevantErrors))

	for idx, err := range relevantErrors {
		reports[idx] = fmt.Sprintf(`*ERROR*
			- *Severity*: %v
			- *Started at*: %s
			- *Workflow*: %v
			- *Message*: %s
			`,
			err.ErrSeverity,
			err.StartedAt.Format(constants.NOTIFICATION_TIME_REP),
			fmt.Sprintf("`%s`", err.workflow),
			err.Err.Error(),
		)
	}

	return reports
}

func isLongRunningError(errStartTime time.Time) bool {
	return time.Since(errStartTime) >= 5*time.Minute
}
