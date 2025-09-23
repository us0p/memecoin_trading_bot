package coinprovider

import (
	"net/http"
	"os"

	"memecoin_trading_bot/app/utils"
)

type paramConfiguration struct {
	Encoding string `json:"encoding"`
}

func newParamConfiguration(encoding string) paramConfiguration {
	return paramConfiguration{
		encoding,
	}
}

type heliusRPCParams struct {
	JsonRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

func (hrp *heliusRPCParams) addParamConfig(config paramConfiguration) {
	hrp.Params = append(hrp.Params, config)
}

func newHeliusRPCParams(method, param string) heliusRPCParams {
	return heliusRPCParams{
		"2.0",
		"1",
		method,
		[]any{
			param,
		},
	}
}

type TokenAuthorities struct {
	// Mint and Freeze authorities are enabled when there's an Account
	// address associated with them.
	// If they're empty, they're not enabled.
	FreezeAuthority string `json:"freezeAuthority"`
	MintAuthority   string `json:"mintAuthority"`
	Mint            string `json:"mint"`
}

func (t TokenAuthorities) GetTokenMint() string {
	return t.Mint
}

type HeliusResponse struct {
	Result struct {
		Value struct {
			Data struct {
				Parsed struct {
					Info TokenAuthorities `json:"info"`
				} `json:"parsed"`
			} `json:"data"`
		} `json:"value"`
	} `json:"result"`
}

func GetTokenAuthorities(client *http.Client, url, mint string) (TokenAuthorities, error) {
	if mint == "" {
		return TokenAuthorities{}, nil
	}
	rpc_params := newHeliusRPCParams(
		"getAccountInfo",
		mint,
	)

	rpc_params.addParamConfig(
		newParamConfiguration("jsonParsed"),
	)

	requester, err := utils.NewRequester[HeliusResponse](client, url, http.MethodPost)
	if err != nil {
		return TokenAuthorities{}, err
	}

	requester.AddHeader("ContentType", "application/json")
	requester.AddQuery("api-key", os.Getenv("HELIUS_API_KEY"))
	req_with_body, err := requester.AddBody(rpc_params)

	if err != nil {
		return TokenAuthorities{}, err
	}

	token_authorities, err := req_with_body.Do()
	if err != nil {
		return TokenAuthorities{}, err
	}

	token_authorities.Result.Value.Data.Parsed.Info.Mint = mint

	return token_authorities.Result.Value.Data.Parsed.Info, nil
}
