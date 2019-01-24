package client

import (
	"io"
	"net/http"
	"path"

	"github.com/opencontainers/go-digest"
)

type Blob interface {
	GetBlob(repo string, digest digest.Digest) (io.Reader, error)
	DeleteBlob(repo string, digest digest.Digest) error
	PutBlob(repo string, digest digest.Digest, blob io.Writer) error
}

// VerifyBlob confirms the existance of a manifest in a remote registry.
func (api *DistributionAPI) VerifyBlob(
	repo string, digest digest.Digest) error {
	c := api.client

	u := *c.Host
	u.Path = path.Join("/v2", repo, "blobs", digest.String())

	req := &http.Request{
		Method: "HEAD",
		URL:    &u,
		Header: make(http.Header),
	}

	resp, err := c.RoundTrip(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrBlobNotExist
	}

	return nil
}
