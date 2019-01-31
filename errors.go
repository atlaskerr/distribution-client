package client

import (
	"errors"
)

// Error satisfies the Error interface.
func (er *ErrorResponse) Error() string {
	return "distribution: registry returned error"
}

func (er *ErrorResponse) Detail() []ErrorInfo {
	return er.Errors
}

type ErrorResponse struct {
	Errors []ErrorInfo `json:"errors"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

var (
	ErrNotOCI           = errors.New("distribution: registry isn't OCI-compliant")
	ErrManifestNotExist = errors.New("distribution: manifest does not exist")
	ErrBlobNotExist     = errors.New("distribution: blob does not exist")
	ErrParseBody        = errors.New("distribtuion: could not read response body")
	ErrParseJSON        = errors.New("distribution: could not parse JSON")
	ErrUnknownMediaType = errors.New("distribution: unknown media type")
	ErrSchemaValidation = errors.New("distribution: unable to validate schema")
	ErrInvalidIndex     = errors.New("distribution: registry returned invalid index")
	ErrInvalidManifest  = errors.New("distribution: registry returned invalid manifest")
)
