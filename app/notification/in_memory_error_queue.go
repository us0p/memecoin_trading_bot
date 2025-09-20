package notification

import (
	"time"
)

type Severity int8

const (
	Transient Severity = iota
	Core
	Fatal
)

type ErrorNotification struct {
	Err         error
	Sent        bool
	StartedAt   time.Time
	ErrSeverity Severity
}

type Workflow string

const (
	PullCoin           Workflow = "PULLCOIN"
	TradeOpEval                 = "TRADE_OP_EVAL"
	TokenAuthorityEval          = "TOKEN_AUTHORITY_EVAL"
	DatabaseOp                  = "DATABASE_OPERATION"
	TokenDataAgg                = "TOKEN_DATA_AGGREGATION"
)

type InMemoryErrorQueueKey struct {
	mint     string
	workflow Workflow
}

type InMemoryErrorQueue map[InMemoryErrorQueueKey][]ErrorNotification

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
