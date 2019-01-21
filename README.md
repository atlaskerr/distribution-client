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

import "log"

func main() {
	cfg := client.Config{
		BaseEndpoint: "http://127.0.0.1:84832",
		Transport: client.DefaultTransport,
		Auth: client.BasicAuth{
			Username: "user",
			Password: "pass",
		},
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	api := c.NewDistributionAPI()
}
```


