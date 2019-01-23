package client

import "errors"

var (

	// ErrNotOCI is an error returned when a remote registry fails an OCI
	// compliance verification check.
	ErrNotOCI = errors.New("distribution: registry isn't OCI-compliant")

	//ErrManifestNotExist is an error returned when a manifest does not exist in
	//a repository.
	ErrManifestNotExist = errors.New("distribution: manifest does not exist")
)
