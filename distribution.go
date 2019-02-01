package client

import (
	"net/url"

	ischema "github.com/atlaskerr/oci-schemas"
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
	imageIndexSchema    *gojsonschema.Schema
	imageManifestSchema *gojsonschema.Schema
}

// NewRegistry returns a fully initialized Registry.
func (api *DistributionAPI) NewRegistry(host string, auth Authenticator) *Registry {
	hostURL, _ := url.Parse(host)

	return &Registry{
		Client *Client
	}
}

// Registry represents a remote OCI-compliant registries.
type Registry struct {
	Client *Client
}
