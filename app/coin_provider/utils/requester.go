package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Requester struct {
	method string
	url    url.URL
	header http.Header
	body   io.ReadCloser
}

func (r *Requester) ConfigHeader(header, value string) {
	r.header.Add(header, value)
}

func (r *Requester) ConfigBody(body string) error {
	var body_as_bytes bytes.Buffer

	decoder := json.NewDecoder(&body_as_bytes)

	if err := decoder.Decode(&body); err != nil {
		return err
	}

	r.body = io.NopCloser(&body_as_bytes)

	return nil
}

func NewRequester(url url.URL, method string) Requester {
	return Requester{
		method: method,
		url:    url,
	}
}
