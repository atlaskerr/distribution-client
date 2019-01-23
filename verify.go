package client

import (
	"errors"
	"net/http"
)

// ErrNotOCI is an error returned when a remote registry fails an OCI compliance
// verification check.
var ErrNotOCI = errors.New("distribution: registry isn't OCI-compliant")

// HeaderVersionCheck is the header defined in the distribution spec that
// clients use to verify that a registry is OCI-compliant.
var HeaderVersionCheck = "Docker-Distribution-Api-Version"

// Verify verifies that the registry is OCI-compliant.
func (api *DistributionAPI) Verify() error {
	c := api.client

	u := *c.Host
	u.Path = "/v2"

	req := &http.Request{
		Method: "GET",
		URL:    &u,
		Header: make(http.Header),
	}
	resp, err := c.RoundTrip(req)
	if err != nil {
		return err
	}

	dockerHeader := resp.Header.Get(HeaderVersionCheck)
	if dockerHeader != "registry/2.0" {
		return ErrNotOCI
	}

	return nil

}
