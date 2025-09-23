package coinprovider

import (
	"memecoin_trading_bot/app/utils"
	"net/http"
	"os"
)

type TransactionSimulation struct {
	Err  string   `json:"err"`
	Logs []string `json:"logs"`
}

type heliusTransactionSimulationReturnData struct {
	Result struct {
		Value TransactionSimulation `json:"value"`
	} `json:"result"`
}

func SimulateTransactionExecution(
	client *http.Client,
	url,
	signedTransaction string,
) (TransactionSimulation, error) {
	requester, err := utils.NewRequester[heliusTransactionSimulationReturnData](
		client,
		url,
		http.MethodPost,
	)
	if err != nil {
		return TransactionSimulation{}, err
	}

	rpc_params := newHeliusRPCParams(
		"simulateTransaction",
		signedTransaction,
	)

	rpc_params.addParamConfig(
		newParamConfiguration("base64"),
	)

	requester.AddHeader("Content-Type", "application/json")
	requester.AddQuery("api-key", os.Getenv("HELIUS_API_KEY"))
	req_with_body, err := requester.AddBody(rpc_params)
	if err != nil {
		return TransactionSimulation{}, err
	}

	transaction_simulation, err := req_with_body.Do()
	if err != nil {
		return TransactionSimulation{}, err
	}

	return transaction_simulation.Result.Value, nil
}
