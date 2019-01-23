package client

import (
	"net/http"
)

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

	dockerHeader := resp.Header.Get(headerVersionCheck)
	if dockerHeader != "registry/2.0" {
		return ErrNotOCI
	}

	return nil

}
