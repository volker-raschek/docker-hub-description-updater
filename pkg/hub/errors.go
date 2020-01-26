package hub

import "errors"

var (
	errorFailedToCreateRequest    = errors.New("Failed to create http request")
	errorFailedToParseJSON        = errors.New("Failed to parse json")
	errorFailedToParseURL         = errors.New("Failed to parse url")
	errorFailedToSendRequest      = errors.New("Failed to send http request")
	errorNoUserDefined            = errors.New("No User defined")
	errorNoPasswordDefined        = errors.New("No Password defined")
	errorNoNamespaceDefined       = errors.New("No Namespace defined")
	errorNoRepositoryDefined      = errors.New("No Repository defined")
	errorUnexpectedHTTPStatuscode = errors.New("Unexpected HTTP-Statuscode")
)
