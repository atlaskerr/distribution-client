package client

import (
	"strings"
)

var (
	// HeaderVersionCheck is the header defined in the distribution spec that
	// clients use to verify that a registry is OCI-compliant.
	headerVersionCheck = "Docker-Distribution-Api-Version"
)

func allowMediaTypes(req *http.Request, mediaTypes ...string) {
	allow := make([]string, len(mediaTypes))
	for _, mt := range mediaTypes {
		allow = append(allow, mt)
	}

	accept := strings.Join(allow, ",")
	req.Header.Set("Accept", accept)
}
