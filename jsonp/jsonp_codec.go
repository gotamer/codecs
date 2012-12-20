package jsonp

import (
	jsonEncoding "encoding/json"
	"errors"
	"github.com/stretchrcom/stretchr/constants"
	"github.com/stretchrcom/stringy"
	"github.com/stretchrcom/web"
)

// ErrorMissingCallback is the error for when a callback option is expected but missing.
var ErrorMissingCallback = errors.New("A callback is required for JSONP")

// ErrorUnmarshalNotSupported is the error for when Unmarshal is called but not supported.
var ErrorUnmarshalNotSupported = errors.New("Unmarshalling an object is not supported for JSONP")

// JsonPCodec converts objects to JSONP.
type JsonPCodec struct{}

// Marshal converts an object to JSONP.
func (c *JsonPCodec) Marshal(object interface{}, options map[string]interface{}) ([]byte, error) {

	if len(options) == 0 {
		return nil, ErrorMissingCallback
	}

	json, err := jsonEncoding.Marshal(object)

	if err != nil {
		return nil, err
	}

	// #codec-context-options
	// the assumption is options[0] is the callback parameter,
	// and options[1] is the client-context (NB: not *Context) string.

	var callbackString string
	var clientContextString string

	clientContextString, hasClientContext := options[constants.OptionKeyClientContext].(string)

	if !hasClientContext {
		callbackString = stringy.Merge(options[constants.OptionKeyCallback].(string), "(", string(json), ");")
	} else {
		callbackString = stringy.Merge(options[constants.OptionKeyCallback].(string), "(", string(json), `,"`, clientContextString, `"`, ");")
	}

	return []byte(callbackString), nil
}

// Unmarshal is not supported for JSONP. Returns an error.
func (c *JsonPCodec) Unmarshal(data []byte, obj interface{}) error {
	return ErrorUnmarshalNotSupported
}

// ContentType returns the content type for this codec.
func (c *JsonPCodec) ContentType() string {
	return web.ContentTypeJavaScript
}

// FileExtensions returns the file extension for this codec.
func (c *JsonPCodec) FileExtensions() string {
	return web.FileExtensionJavaScript
}

// CanMarshalWithCallback returns whether this codec is capable of marshalling a response containing a callback.
func (c *JsonPCodec) CanMarshalWithCallback() bool {
	return true
}