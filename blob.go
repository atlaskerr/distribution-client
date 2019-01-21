package client

import (
	"io"

	"github.com/opencontainers/go-digest"
)

type Blob interface {
	VerifyBlob(repo string, digest digest.Digest) (bool, error)
	GetBlob(repo string, digest digest.Digest) (io.Reader, error)
	DeleteBlob(repo string, digest digest.Digest) error
	PutBlob(repo string, digest digest.Digest, blob io.Reader) error
}
