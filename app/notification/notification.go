package notification

import (
	"context"
	"errors"
	"log"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/entities"
	"memecoin_trading_bot/app/utils"
	"net/http"
	"os"
	"sync"
)

type Notifications struct {
	ErrQueue     InMemoryErrorQueue
	TradesOpened []TradeOpening
	TradesClosed []TradeClosing
	mut          *sync.RWMutex
}

func NewNotificationState() Notifications {
	return Notifications{
		ErrQueue:     make(InMemoryErrorQueue),
		TradesOpened: make([]TradeOpening, 0),
		TradesClosed: make([]TradeClosing, 0),
		mut:          &sync.RWMutex{},
	}
}

func (n *Notifications) RecordError(token string, workflow Workflow, err error, sev Severity) {
	key := newInMemoryErrorQueueKey(token, workflow)

	queue := (*n).ErrQueue[key]

	for _, errNotification := range queue {
		if errors.Is(errNotification.Err, err) {
			return
		}
	}

	n.mut.Lock()
	defer n.mut.Unlock()
	(*n).ErrQueue[key] = append(queue, newErrorNotification(err, sev))
}

func (n *Notifications) RecordTradeExecution(
	db_client *db.DB,
	mint string,
	operation entities.Operation,
	total_sol_wallet float64,
) {
	ctx := context.Background()
	trade_buy, err := db_client.GetTradeNotificationData(ctx, mint, entities.BUY)
	if err != nil {
		n.RecordError(
			mint,
			TradeNotification,
			err,
			Fatal,
		)
		return
	}
	if operation == entities.BUY {
		n.recordTradeOpening(trade_buy, total_sol_wallet)
		return
	} else {
		trade_sell, err := db_client.GetTradeNotificationData(ctx, mint, entities.SELL)
		if err != nil {
			n.RecordError(
				mint,
				TradeNotification,
				err,
				Fatal,
			)
			return
		}
		n.recordTradeClosing(trade_buy, trade_sell)
	}
}

func (n *Notifications) SendNotifications(client *http.Client, telegram_url string) {
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

func (n *Notifications) notifyTelegram(client *http.Client, telegram_url, message string) error {
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
