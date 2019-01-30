package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	ischema "github.com/atlaskerr/oci-schemas"
	ispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/xeipuuv/gojsonschema"
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

// GetManifests returns an image index and a slice of manifests from a remote
// registry. If the manifest is an Image Index, all manifests referenced in the
// index will be downloaded and returned. If a nil Index is returned, there will
// only be one manifest in the slice.
func (api *DistributionAPI) GetManifests(
	repo string, reference string) (*ispec.Index, *[]ispec.Manifest, error) {
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
		return nil, nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, ErrParseBody
	}

	loader := gojsonschema.NewBytesLoader(b)

	mediaType := resp.Header.Get("Content-Type")
	switch mediaType {
	case ispec.MediaTypeImageIndex:
		schema := ischema.ImageIndexSchema()
		result, err := schema.Validate(loader)
		if err != nil {
			return nil, nil, ErrSchemaValidation
		}

		if !result.Valid() {
			return nil, nil, ErrInvalidIndex
		}
		var idx ispec.Index
		if err := json.Unmarshal(b, &idx); err != nil {
			return nil, nil, ErrParseJSON
		}
		return nil, api.getManifests(idx)
	case ispec.MediaTypeImageManifest:
		var m ispec.Manifest
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, nil, ErrParseJSON
		}
		return nil, &[]ispec.Manifest{m}, nil
	}
	return nil, nil, ErrUnknownMediaType
}

func (api *DistributionAPI) getManifests(idx ispec.Index) (*[]ispec.Manifest, error) {
	var m []ispec.Manifest
	return &m, nil
}
