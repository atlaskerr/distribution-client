package client

import "errors"

var (

	// ErrNotOCI is an error returned when a remote registry fails an OCI
	// compliance verification check.
	ErrNotOCI = errors.New("distribution: registry isn't OCI-compliant")

	// ErrManifestNotExist is an error returned when a manifest does not exist in
	// a repository.
	ErrManifestNotExist = errors.New("distribution: manifest does not exist")

	// ErrBlobNotExist is an error returned when a blob does not exist in
	// a repository.
	ErrBlobNotExist = errors.New("distribution: blob does not exist")

	ErrParseBody = errors.New("distribtuion: could not read response body")

	ErrParseJSON = errors.New("distribution: could not parse JSON")

	ErrUnknownMediaType = errors.New("distribution: unknown media type")

	ErrSchemaValidation = errors.New("distribution: unable to validate schema")

	ErrInvalidIndex = errors.New("distribution: invalid image index")
)
