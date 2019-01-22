package client

import (
	"errors"
	"net/http"
)

// ErrNotOCI is an error returned when a remote registry fails an OCI compliance
// verification check.
var ErrNotOCI = errors.New("distribution: registry isn't OCI-compliant")

// Verify verifies that the registry is OCI-compliant.
func (api *DistributionAPI) Verify() error {
	c := api.client

	u := *c.Host
	u.Path = "/v2"

	req := &http.Request{
		Method: "GET",
		URL:    &u,
	}
	resp, err := c.RoundTrip(req)
	if err != nil {
		return err
	}

	dockerHeader := resp.Header.Get("Docker-Distribution-Api-Version")
	if dockerHeader != "registry/2.0" {
		return ErrNotOCI
	}

	return nil

}
