package client

import (
	"net/http"
)

// Verify verifies that the registry is OCI-compliant.
func (api *DistributionAPI) Verify() (bool, error) {
	c := api.client

	u := *c.BaseEndpoint
	u.Path = "/v2"

	req := &http.Request{
		Method: "GET",
		URL:    &u,
	}
	resp, err := c.RoundTrip(req)
	if err != nil {
		return false, err
	}

	dockerHeader := resp.Header.Get("Docker-Distribution-Api-Version")
	if dockerHeader != "registry/2.0" {
		return false, nil
	}

	return true, nil

}
