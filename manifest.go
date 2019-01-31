package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	digest "github.com/opencontainers/go-digest"
	ispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/xeipuuv/gojsonschema"
)

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

	mediaType := resp.Header.Get("Content-Type")

	if mediaType == ispec.MediaTypeImageIndex {
		idx, err := api.parseIndex(resp.Body)
		if err != nil {
			return nil, nil, err
		}

		for _, desc := range idx.Manifests {
			_, _ = api.getManifest(repo, desc.Digest)
		}

	}

	return nil, nil, ErrUnknownMediaType
}

func (api *DistributionAPI) getManifest(repo string, digest digest.Digest) (*ispec.Manifest, error) {
	c := api.client

	u := *c.Host
	u.Path = path.Join("/v2", repo, "manifests", digest.String())

	req := &http.Request{
		Method: "GET",
		URL:    &u,
		Header: make(http.Header),
	}
	resp, err := c.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrParseBody
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp *ErrorResponse
		if err := json.Unmarshal(b, errorResp); err != nil {
			return nil, ErrParseJSON
		}
		return nil, errorResp
	}

	var manifest *ispec.Manifest
	if err := json.Unmarshal(b, manifest); err != nil {
		return nil, ErrParseJSON
	}
	return manifest, nil
}

func (api *DistributionAPI) parseIndex(idx io.Reader) (*ispec.Index, error) {
	b, err := ioutil.ReadAll(idx)
	if err != nil {
		return nil, err
	}

	loader := gojsonschema.NewBytesLoader(b)
	res, err := api.imageIndexSchema.Validate(loader)
	if err != nil {
		return nil, ErrSchemaValidation
	}

	if !res.Valid() {
		return nil, ErrInvalidIndex
	}

	var index *ispec.Index
	err = json.Unmarshal(b, index)
	if err != nil {
		return nil, ErrParseJSON
	}

	return index, nil
}
