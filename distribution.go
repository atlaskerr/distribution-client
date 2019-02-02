package client

import (
	"net/http"
	"net/url"

	ischema "github.com/atlaskerr/oci-schemas"
	"github.com/xeipuuv/gojsonschema"
)

// NewDistributionAPI returns a fully initialized API for interacting with
// a remote registry.
func NewDistributionAPI(c *Client) *DistributionAPI {
	api := &DistributionAPI{
		client:              c,
		imageIndexSchema:    ischema.ImageIndexSchema(),
		imageManifestSchema: ischema.ImageManifestSchema(),
	}
	return api
}

// DistributionAPI contains methods for interacting with a remote registry.
type DistributionAPI struct {
	client              *Client
	validator           *Validator
	imageIndexSchema    *gojsonschema.Schema
	imageManifestSchema *gojsonschema.Schema
}

// NewRegistry returns a fully initialized Registry.
func (api *DistributionAPI) NewRegistry(host string, auth *Authenticator) *Registry {
	hostURL, _ := url.Parse(host)

	return &Registry{
		client: api.client,
		Host:   hostURL,
		Auth:   *auth,
	}
}

// Registry represents a remote OCI-compliant registries.
type Registry struct {
	client *Client
	Host   *url.URL
	Auth   Authenticator
}

// RoundTrip is the Registry implementation of http.RoundTripper.
func (r *Registry) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if r.Auth != nil {
		r.Auth.Set(req)
	}
	return r.client.RoundTrip(req)
}
