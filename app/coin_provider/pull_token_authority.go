package coinprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"memecoin_trading_bot/app/app_errors"
	"net/http"
	"os"
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

func newHeliusRPCParams(method string, mint string) heliusRPCParams {
	return heliusRPCParams{
		"2.0",
		"1",
		method,
		[]any{
			mint,
		},
	}
}

type TokenAuthorities struct {
	// Mint and Freeze authorities are enabled when there's an Account
	// address associated with them.
	// If they're empty, they're not enabled.
	FreezeAuthority string `json:"freezeAuthority"`
	MintAuthority   string `json:"mintAuthority"`
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
	rpc_params := newHeliusRPCParams(
		"getAccountInfo",
		mint,
	)

	rpc_params.addParamConfig(
		newParamConfiguration("jsonParsed"),
	)

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(rpc_params); err != nil {
		return TokenAuthorities{}, err
	}

	api_key := "?api-key=" + os.Getenv("HELIUS_API_KEY")
	url_with_key := url + api_key

	res, err := client.Post(url_with_key, "application/json", &buf)
	if err != nil {
		return TokenAuthorities{}, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		error_response, err := io.ReadAll(res.Body)
		if err != nil {
			return TokenAuthorities{}, fmt.Errorf("%w, %s",
				app_errors.ErrReadingAPIResponse, err,
			)
		}

		return TokenAuthorities{}, fmt.Errorf(
			"%w, Status: %d while calling Helius 'getAccountInfo'. Error: %s",
			app_errors.ErrNonOkStatus,
			res.StatusCode,
			error_response,
		)
	}

	var token_authorities HeliusResponse
	if err = decoder.Decode(&token_authorities); err != nil {
		return TokenAuthorities{}, err
	}

	return token_authorities.Result.Value.Data.Parsed.Info, nil
}
