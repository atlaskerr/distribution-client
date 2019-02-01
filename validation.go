package client

import (
	"github.com/xeipuuv/gojsonschema"
)

type Validator interface {
	ValidateImageIndex()
}
