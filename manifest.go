package client

import (
	"errors"
	"net/http"
	"path"

	ispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type Manifest interface {
	GetManifest(repo, reference string) (ispec.Manifest, error)
	DeleteManifest(repo, reference string) error
	PutManifest(repo, reference string, manifest ispec.Manifest) error
}

var ErrManifestNotExist = errors.New("distribution: manifest does not exist")

// VerifyManifest confirms the existance of a manifest in a remote registry.
func (api *DistributionAPI) VerifyManifest(repo, reference string) error {
	c := api.client

	u := *c.Host
	u.Path = path.Join("/v2", repo, "manifests", reference)

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
		return ErrManifestNotExist
	}

	// TODO (atlaskerr): Move block to GetManifest method.
	//manifestSlice := []string{
	//	ispec.MediaTypeImageManifest,
	//	ispec.MediaTypeImageIndex,
	//}
	//acceptHeaderValue := strings.Join(manifestSlice, ",")
	//req.Header.Set("Accept", acceptHeaderValue)

	return nil
}
