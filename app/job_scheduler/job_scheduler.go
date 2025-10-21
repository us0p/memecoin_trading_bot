package jobscheduler

import (
	"log"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/notification"
	"memecoin_trading_bot/app/workflows"
	"net/http"
	"time"
)

type Workflow func(
	*http.Client,
	*db.DB,
	*notification.Notifications,
	*workflows.TransactionProcessing,
)

type Job struct {
	Workflow     Workflow
	WorkflowName string
}

type Interval time.Duration

const (
	FIVE_SECOND_INTERVAL Interval = Interval(time.Second * 5)
	ONE_MINUTE_INTERVAL  Interval = Interval(time.Minute)
	FIVE_MINUTE_INTERVAL Interval = Interval(time.Minute * 5)
)

type JobSchedulerMap map[Interval][]Job

func NewJobSchedulerMap() JobSchedulerMap {
	return make(JobSchedulerMap)
}

func (j *JobSchedulerMap) RegisterJob(workflow Workflow, workflowName string, interval Interval) {
	(*j)[interval] = append((*j)[interval], Job{workflow, workflowName})
}

func (j *JobSchedulerMap) jobExecutor(
	workflow Workflow,
	http_client *http.Client,
	db_client *db.DB,
	tp *workflows.TransactionProcessing,
	nf_state *notification.Notifications,
) {
	workflow(http_client, db_client, nf_state, tp)

	nf_state.SendNotifications(
		http_client,
		constants.TELEGRAM_API_URL,
	)
}

func (j *JobSchedulerMap) jobsExecutor(
	http_client *http.Client,
	db_client *db.DB,
	nf_state *notification.Notifications,
	tp *workflows.TransactionProcessing,
	interval Interval,
) {
	for {
		for _, job := range (*j)[interval] {
			log.Printf("Executing workflow: %s...\n", job.WorkflowName)
			go j.jobExecutor(
				job.Workflow,
				http_client,
				db_client,
				tp,
				nf_state,
			)
		}
		time.Sleep(time.Duration(interval))
	}
}

func (j *JobSchedulerMap) StartJobExecutor(
	http_client *http.Client,
	db_client *db.DB,
	nf_state *notification.Notifications,
	tp *workflows.TransactionProcessing,
) {
	for interval := range *j {
		go j.jobsExecutor(
			http_client,
			db_client,
			nf_state,
			tp,
			interval,
		)
	}
}
