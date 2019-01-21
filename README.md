# distribution-client

distribution-client is the Go client library for the OCI Distribution
Specification.

## Install
```bash
go get github.com/atlaskerr/distribution-client
```

## Usage

```go
package main

import (
	"log"
	"time"
)

func main() {
	cfg := client.Config{
		Endpoint: "http://127.0.0.1:84832",
		Transport: client.DefaultTransport,

		// Set timeout per request to fail fast when the registry is
		// unavailable.
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	api := c.NewDistributionAPI()
}
```


