package hub

import "errors"

var (
	errFailedToCreateRequest    = errors.New("failed to create http request")
	errFailedToParseJSON        = errors.New("failed to parse json")
	errFailedToParseURL         = errors.New("failed to parse url")
	errFailedToSendRequest      = errors.New("failed to send http request")
	errNoNamespaceDefined       = errors.New("no Namespace defined")
	errNoRepositoryDefined      = errors.New("no Repository defined")
	errUnexpectedHTTPStatuscode = errors.New("unexpected HTTP-Statuscode")
)
