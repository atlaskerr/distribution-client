package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	ispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// VerifyManifest confirms the existance of a manifest in a remote registry.
func (r *Registry) VerifyManifest(img Image) error {
	c := r.client

	req := new(http.Request)
	req.Method = "HEAD"
	req.URL = r.Host
	req.URL.Path = manifestEndpoint(img.Repository, img.Reference)
	req.Header = make(http.Header)

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
func (r *Registry) GetManifests(img Image) (*ispec.Index, *[]ispec.Manifest, error) {
	c := r.client

	req := new(http.Request)
	req.Method = "GET"
	req.URL = r.Host
	req.URL.Path = manifestEndpoint(img.Repository, img.Reference)
	req.Header = make(http.Header)

	req.Header["Accept"] = []string{
		ispec.MediaTypeImageIndex,
		ispec.MediaTypeImageManifest,
	}

	resp, err := c.RoundTrip(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var idx *ispec.Index
	manifests := make([]ispec.Manifest, 0)

	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case ispec.MediaTypeImageIndex:
		err := validateIndex(resp.Body)
		if err != nil {
			return nil, nil, err
		}
		parseIndex(resp.Body, idx)

	case ispec.MediaTypeImageManifest:
		err := validateManifest(resp.Body)
		if err != nil {
			return nil, nil, err
		}
		var m *ispec.Manifest
		parseManifest(resp.Body, m)
		manifests = append(manifests, *m)

	default:
		return nil, nil, ErrUnknownMediaType
	}

	return idx, &manifests, nil
}

func (r *Registry) getManifest(img Image) (*ispec.Manifest, error) {
	c := r.client

	req := new(http.Request)
	req.Method = "GET"
	req.URL = r.Host
	req.URL.Path = manifestEndpoint(img.Repository, img.Reference)
	req.Header = make(http.Header)

	req.Header["Accept"] = []string{ispec.MediaTypeImageManifest}

	resp, err := c.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp *ErrorResponse
		if err := json.Unmarshal(b, errorResp); err != nil {
			return nil, ErrParseJSON
		}
		return nil, errorResp
	}

	var m *ispec.Manifest
	if err := parseManifest(resp.Body, m); err != nil {
		return nil, err
	}

	return m, nil
}

func parseManifest(data io.Reader, manifest *ispec.Manifest) error {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, manifest); err != nil {
		return ErrParseJSON
	}
	return nil
}

func parseIndex(data io.Reader, idx *ispec.Index) error {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, idx); err != nil {
		return ErrParseJSON
	}

	return nil
}

func manifestEndpoint(repo, reference string) string {
	return path.Join("/v2", repo, "manifests", reference)
}
