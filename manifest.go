package client

import (
	ispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type Manifest interface {
	VerifyManifest(repo, reference string) (bool, error)
	GetManifest(repo, reference string) (ispec.Manifest, error)
	DeleteManifest(repo, reference string) error
	PutManifest(repo, reference string, manifest ispec.Manifest) error
}
