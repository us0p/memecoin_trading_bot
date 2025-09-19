package notification

import (
	"errors"
	"fmt"
	"log"
	"memecoin_trading_bot/app/coin_provider/utils"
	"net/http"
	"os"
	"sync"
	"time"
)

func NewNotificationState() Notifications {
	return Notifications{
		ErrQueue:     make(InMemoryErrorQueue),
		TradesOpened: make([]TradeOpening, 0),
		TradesClosed: make([]TradeClosing, 0),
	}
}

func (n Notifications) notifyTelegram(client *http.Client, telegram_url, message string) error {
	url_with_token := telegram_url + os.Getenv("TELEGRAM_TOKEN")
	requester, err := utils.NewRequester[any](client, url_with_token, http.MethodGet)
	if err != nil {
		return err
	}

	requester.AddPath("/sendMessage")
	requester.AddQuery("parse_mode", "Markdown")
	requester.AddQuery("chat_id", os.Getenv("TELEGRAM_CHAT_ID"))
	requester.AddQuery("text", message)

	_, err = requester.Do()
	if err != nil {
		return err
	}

	return nil
}

func (n *Notifications) RecordError(token string, workflow Workflow, err error, sev Severity) {
	key := newInMemoryErrorQueueKey(token, workflow)

	queue := (*n).ErrQueue[key]

	for _, errNotification := range queue {
		if errors.Is(errNotification.Err, err) {
			return
		}
	}

	(*n).ErrQueue[key] = append(queue, newErrorNotification(err, sev))
}

func (n Notifications) SendNotifications(client *http.Client, telegram_url string) {
	reports := [][]string{
		n.errReport(),
		n.openTradeReport(),
		n.closeTradeReport(),
	}

	for _, report_queue := range reports {
		var wg sync.WaitGroup
		for _, report := range report_queue {
			wg.Add(1)
			go func(report string) {
				defer wg.Done()
				err := n.notifyTelegram(client, telegram_url, report)
				if err != nil {
					log.Printf("Error while sending telegram message. ERROR: %s\n", err)
				}
			}(report)
		}
		wg.Wait()
	}
}

// New line encoded for URL representation
const nl_url_ecoded = "%0A"
const time_rep = "02-01-2006 03:04PM"

func (n Notifications) openTradeReport() []string {
	reports := make([]string, len(n.TradesOpened))

	for idx, t := range n.TradesOpened {
		reports[idx] = fmt.Sprintf(`*TRADE OPENING*
			- *Symbol*: %s
			- *Opened at*: %s
			- *Amount SOL*: %f
			- *Wallet percentage*: %f%%`,
			t.Symbol,
			t.OpenedAt.Format(time_rep),
			t.SolAmount,
			t.WalletPercentage,
		)
	}

	return reports
}

func (n Notifications) closeTradeReport() []string {
	reports := make([]string, len(n.TradesClosed))

	for idx, t := range n.TradesClosed {
		reports[idx] = fmt.Sprintf(`*TRADE CLOSED*
			- *Symbol*: %s
			- *Closed at*: %s
			- *Duration*: %s
			- *Profit %%*: %f%%
			- *Pofit in SOL*: %f
			- *Operation Status:* %s`,
			t.Symbol,
			t.ClosedAt.Format(time_rep),
			t.Duration,
			t.ProfitPercentage,
			t.ProfitSOL,
			t.OpStatus,
		)
	}

	return reports
}

type ErrorReport struct {
	InMemoryErrorQueueKey
	ErrorNotification
}

func (n Notifications) getRelevantErrors() []ErrorReport {
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

func (n Notifications) errReport() []string {
	relevantErrors := n.getRelevantErrors()

	reports := make([]string, len(relevantErrors))

	for idx, err := range relevantErrors {
		reports[idx] = fmt.Sprintf(`*ERROR*
			- *Severity*: %v
			- *Started at*: %s
			- *Message*: %s
			`,
			err.ErrSeverity,
			err.StartedAt.Format(time_rep),
			err.Err.Error(),
		)
	}

	return reports
}

func isLongRunningError(errStartTime time.Time) bool {
	return time.Since(errStartTime) >= 5*time.Minute
}
