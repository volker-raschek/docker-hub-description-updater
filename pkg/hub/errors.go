package hub

import "errors"

var (
	errorNoUserDefined       = errors.New("No User defined")
	errorNoPasswordDefined   = errors.New("No Password defined")
	errorNoNamespaceDefined  = errors.New("No Namespace defined")
	errorNoRepositoryDefined = errors.New("No Repository defined")
)
