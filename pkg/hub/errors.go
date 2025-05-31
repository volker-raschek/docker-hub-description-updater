package hub

import "errors"

var (
	errorFailedToCreateRequest    = errors.New("failed to create http request")
	errorFailedToParseJSON        = errors.New("failed to parse json")
	errorFailedToParseURL         = errors.New("failed to parse url")
	errorFailedToSendRequest      = errors.New("failed to send http request")
	errorNoNamespaceDefined       = errors.New("no Namespace defined")
	errorNoRepositoryDefined      = errors.New("no Repository defined")
	errorUnexpectedHTTPStatuscode = errors.New("unexpected HTTP-Statuscode")
)
