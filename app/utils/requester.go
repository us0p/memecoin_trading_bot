package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"memecoin_trading_bot/app/app_errors"
)

type Requester[T any] struct {
	client *http.Client
	method string
	url    url.URL
	header http.Header
	body   io.ReadCloser
}

func (r *Requester[T]) AddHeader(header, value string) *Requester[T] {
	r.header.Add(header, value)
	return r
}

func (r *Requester[T]) AddBody(body any) (*Requester[T], error) {
	var body_as_bytes bytes.Buffer

	if err := json.NewEncoder(&body_as_bytes).Encode(body); err != nil {
		return &Requester[T]{}, err
	}

	r.body = io.NopCloser(&body_as_bytes)

	return r, nil
}

func (r *Requester[T]) AddQuery(key, value string) *Requester[T] {
	q := r.url.Query()
	q.Add(key, value)

	r.url.RawQuery = q.Encode()
	return r
}

func (r *Requester[T]) AddPath(path string) *Requester[T] {
	r.url = *r.url.JoinPath(path)
	return r
}

func (r *Requester[T]) buildRequest() http.Request {
	return http.Request{
		Method: r.method,
		URL:    &r.url,
		Header: r.header,
		Body:   r.body,
	}
}

func (r *Requester[T]) Do() (T, error) {
	req := r.buildRequest()

	res, err := r.client.Do(&req)

	// Represents the zero value for the dynamic type.
	// If T is a struct it'll work the same by zeroing all the fields of the struct.
	var responseBody T
	if err != nil {
		return responseBody, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		error_response, err := io.ReadAll(res.Body)
		if err != nil {
			return responseBody, fmt.Errorf("%w, %s",
				app_errors.ErrReadingAPIResponse, err,
			)
		}

		return responseBody, fmt.Errorf(
			"%w, URL: %s, Status: %d, Body: %s",
			app_errors.ErrNonOkStatus,
			r.url.String(),
			res.StatusCode,
			error_response,
		)
	}

	if err = decoder.Decode(&responseBody); err != nil {
		return responseBody, err
	}

	return responseBody, nil
}

func NewRequester[T any](client *http.Client, url_raw, method string) (Requester[T], error) {
	parsed_url, err := url.Parse(url_raw)

	if err != nil {
		return Requester[T]{}, err
	}

	return Requester[T]{
		client: client,
		method: method,
		url:    *parsed_url,
		header: make(http.Header),
	}, nil
}
