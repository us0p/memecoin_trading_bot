package notification

import "time"

type Severity int8

const (
	Fatal     Severity = 2
	Core      Severity = 1
	Transient Severity = 0
)

type ErrorNotification struct {
	Err         error
	Sent        bool
	StartedAt   time.Time
	ErrSeverity Severity
}

type Workflow string

const (
	PullCoin Workflow = "PULLCOIN"
)

type InMemoryErrorQueueKey struct {
	mint     string
	workflow Workflow
}

type InMemoryErrorQueue map[InMemoryErrorQueueKey][]ErrorNotification

type TradeOpening struct {
	Symbol           string
	OpenedAt         time.Time
	SolAmount        float64
	WalletPercentage float64
}

type Status string

const (
	Win  Status = "WIN"
	Loss Status = "LOSS"
)

type TradeClosing struct {
	Symbol           string
	ClosedAt         time.Time
	Duration         string
	ProfitPercentage float32
	ProfitSOL        float64
	OpStatus         Status
}

type Notifications struct {
	ErrQueue     InMemoryErrorQueue
	TradesOpened []TradeOpening
	TradesClosed []TradeClosing
}
