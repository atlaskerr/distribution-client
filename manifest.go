package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	ispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type Manifest interface {
	DeleteManifest(repo, reference string) error
	PutManifest(repo, reference string, manifest ispec.Manifest) error
}

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

	return nil
}

// GetManifest retrieves a manifest from a remote registry.
func (api *DistributionAPI) GetManifest(
	repo string, reference string) (*ispec.Manifest, error) {
	c := api.client

	u := *c.Host
	u.Path = path.Join("/v2", repo, "manifests", reference)

	req := &http.Request{
		Method: "GET",
		URL:    &u,
		Header: make(http.Header),
	}

	manifestSlice := []string{
		ispec.MediaTypeImageManifest,
		ispec.MediaTypeImageIndex,
	}

	acceptHeaderValue := strings.Join(manifestSlice, ",")
	req.Header.Set(headerAccept, acceptHeaderValue)

	resp, err := c.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	var m *ispec.Manifest
	if err := json.Unmarshal(b, m); err != nil {
		return nil, err
	}
	return m, nil
}
