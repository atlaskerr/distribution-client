package client

import (
	"io"
	"io/ioutil"

	ischema "github.com/atlaskerr/oci-schemas"
	"github.com/xeipuuv/gojsonschema"
)

var (
	imageIndexSchema    = ischema.ImageIndexSchema()
	imageManifestSchema = ischema.ImageManifestSchema()
)

func validateIndex(data io.Reader) error {
	return validate(data, *imageIndexSchema)
}

func validateManifest(data io.Reader) error {
	return validate(data, *imageManifestSchema)
}

func validate(data io.Reader, schema gojsonschema.Schema) error {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	loader := gojsonschema.NewBytesLoader(b)
	res, err := schema.Validate(loader)
	if err != nil {
		return ErrSchemaValidation
	}

	if !res.Valid() {
		e := res.Errors()
		verr := ValidationError(e)
		return &verr
	}

	return nil
}
