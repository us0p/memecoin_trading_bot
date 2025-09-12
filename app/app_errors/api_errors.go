package app_errors

import "errors"

var (
	ErrReadingAPIResponse = errors.New("Failed to read message from API")
	ErrNonOkStatus        = errors.New("Received non-ok status from API")
)
