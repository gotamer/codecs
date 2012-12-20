package services

import (
	"github.com/stretchrcom/codecs"
	"github.com/stretchrcom/codecs/test"
	"github.com/stretchrcom/objects"
	"github.com/stretchrcom/testify/assert"
	"github.com/stretchrcom/testify/mock"
	"github.com/stretchrcom/web"
	"strings"
	"testing"
)

/*
	Test code
*/

func TestInterface(t *testing.T) {
	assert.Implements(t, (*CodecService)(nil), new(WebCodecService), "WebCodecService")
}

func TestInstalledCodecs(t *testing.T) {
	assert.NotNil(t, InstalledCodecs, "InstalledCodecs")
	assert.Equal(t, len(InstalledCodecs), 3, "Should be three codecs installed.")
}

func TestGetCodecForRequest(t *testing.T) {

	service := new(WebCodecService)
	var codec codecs.Codec

	codec, _ = service.GetCodecForRequest(web.ContentTypeJson)

	if assert.NotNil(t, codec, "Json should exist") {
		assert.Equal(t, web.ContentTypeJson, codec.ContentType(), "ContentTypeJson")
	}

	// case insensitivity
	codec, _ = service.GetCodecForRequest(strings.ToUpper(web.ContentTypeJson))

	if assert.NotNil(t, codec, "Content case should not matter") {
		assert.Equal(t, web.ContentTypeJson, codec.ContentType(), "ContentTypeJson")
	}

	// default
	codec, _ = service.GetCodecForRequest("")

	if assert.NotNil(t, codec, "Empty contentType string should assume JSON") {
		assert.Equal(t, web.ContentTypeJson, codec.ContentType(), "Should assume JSON.")
	}

}

func TestGetCodecForResponding_DefaultCodec(t *testing.T) {

	service := new(WebCodecService)
	var codec codecs.Codec

	codec, _ = service.GetCodecForResponding("", "", false)

	if assert.NotNil(t, codec, "Return of GetCodecForAcceptStringOrExtension should default to JSON") {
		assert.Equal(t, web.ContentTypeJson, codec.ContentType(), "Should default to JSON")
	}

}

func TestGetCodecForResponding(t *testing.T) {

	service := new(WebCodecService)
	var codec codecs.Codec

	// JSON - accept header

	codec, _ = service.GetCodecForResponding("something/something,application/json,text/xml", "", false)

	if assert.NotNil(t, codec, "Return of GetCodecForAcceptStringOrExtension") {
		assert.Equal(t, web.ContentTypeJson, codec.ContentType(), "ContentTypeJson 1")
	}

	// JSON - accept header (case)

	codec, _ = service.GetCodecForResponding("something/something,application/JSON,text/xml", "", false)

	if assert.NotNil(t, codec, "Case should not matter") {
		assert.Equal(t, web.ContentTypeJson, codec.ContentType(), "Case should not matter")
	}

	// JSON - file extension

	codec, _ = service.GetCodecForResponding("", web.FileExtensionJson, false)

	if assert.NotNil(t, codec, "Return of GetCodecForAcceptStringOrExtension") {
		assert.Equal(t, web.ContentTypeJson, codec.ContentType(), "ContentTypeJson")
	}

	// JSONP - has callback

	codec, _ = service.GetCodecForResponding("", "", true)

	if assert.NotNil(t, codec, "Should return the first codec that can handle a callback") {
		assert.Equal(t, web.ContentTypeJavaScript, codec.ContentType(), "ContentTypeJavaScript")
	}

	// JSONP - file extension

	codec, _ = service.GetCodecForResponding("", web.FileExtensionJavaScript, false)

	if assert.NotNil(t, codec, "Return of GetCodecForAcceptStringOrExtension") {
		assert.Equal(t, web.ContentTypeJavaScript, codec.ContentType(), "ContentTypeJavaScript")
	}

	// JSONP - file extension (case)

	codec, _ = service.GetCodecForResponding("", strings.ToUpper(web.FileExtensionJavaScript), false)

	if assert.NotNil(t, codec, "Return of GetCodecForAcceptStringOrExtension") {
		assert.Equal(t, web.ContentTypeJavaScript, codec.ContentType(), "ContentTypeJavaScript 4")
	}

	// JSONP - Accept header

	codec, _ = service.GetCodecForResponding("something/something,text/javascript,text/xml", "", false)

	if assert.NotNil(t, codec, "Return of GetCodecForAcceptStringOrExtension") {
		assert.Equal(t, web.ContentTypeJavaScript, codec.ContentType(), "ContentTypeJavaScript 5")
	}

}

func TestMarshalWithCodec(t *testing.T) {

	testCodec := new(test.TestCodec)
	service := new(WebCodecService)

	// make some test stuff
	var bytesToReturn []byte = []byte("Hello World")
	var object objects.Map = objects.Map{"Name": "Mat"}
	var option1 string = "Option One"
	var option2 string = "Option Two"

	args := map[string]interface{}{option1: option1, option2: option2}

	// setup expectations
	testCodec.On("Marshal", object, args).Return(bytesToReturn, nil)

	bytes, err := service.MarshalWithCodec(testCodec, object, args)

	if assert.Nil(t, err) {
		assert.Equal(t, string(bytesToReturn), string(bytes))
	}

	// assert that our expectations were met
	mock.AssertExpectationsForObjects(t, testCodec.Mock)

}

func TestMarshalWithCodec_WithFacade(t *testing.T) {

	// func (s *WebCodecService) MarshalWithCodec(codec codecs.Codec, object interface{}, options ...interface{}) ([]byte, error) {

	testCodec := new(test.TestCodec)
	service := new(WebCodecService)

	// make some test stuff
	var bytesToReturn []byte = []byte("Hello World")
	testObjectWithFacade := new(test.TestObjectWithFacade)
	object := objects.Map{"Name": "Mat"}
	var option1 string = "Option One"
	var option2 string = "Option Two"

	args := map[string]interface{}{option1: option1, option2: option2}

	// setup expectations
	testObjectWithFacade.On("PublicData", args).Return(object, nil)
	testCodec.On("Marshal", object, args).Return(bytesToReturn, nil)

	bytes, err := service.MarshalWithCodec(testCodec, testObjectWithFacade, args)

	if assert.Nil(t, err) {
		assert.Equal(t, string(bytesToReturn), string(bytes))
	}

	// assert that our expectations were met
	mock.AssertExpectationsForObjects(t, testCodec.Mock, testObjectWithFacade.Mock)

}

func TestMarshalWithCodec_WithFacade_AndError(t *testing.T) {

	// func (s *WebCodecService) MarshalWithCodec(codec codecs.Codec, object interface{}, options ...interface{}) ([]byte, error) {

	testCodec := new(test.TestCodec)
	service := new(WebCodecService)

	// make some test stuff
	testObjectWithFacade := new(test.TestObjectWithFacade)
	var option1 string = "Option One"
	var option2 string = "Option Two"

	args := map[string]interface{}{option1: option1, option2: option2}

	// setup expectations
	testObjectWithFacade.On("PublicData", args).Return(nil, assert.AnError)

	_, err := service.MarshalWithCodec(testCodec, testObjectWithFacade, args)

	assert.Equal(t, assert.AnError, err)

}

func TestMarshalWithCodec_WithError(t *testing.T) {

	// func (s *WebCodecService) MarshalWithCodec(codec codecs.Codec, object interface{}, options ...interface{}) ([]byte, error) {

	testCodec := new(test.TestCodec)
	service := new(WebCodecService)

	// make some test stuff
	object := objects.Map{"Name": "Mat"}
	var option1 string = "Option One"
	var option2 string = "Option Two"

	args := map[string]interface{}{option1: option1, option2: option2}

	// setup expectations
	testCodec.On("Marshal", object, args).Return(nil, assert.AnError)

	_, err := service.MarshalWithCodec(testCodec, object, args)

	assert.Equal(t, assert.AnError, err, "The error should get returned")

	// assert that our expectations were met
	mock.AssertExpectationsForObjects(t, testCodec.Mock)

}

func TestUnmarshalWithCodec(t *testing.T) {

	// func (s *WebCodecService) UnmarshalWithCodec(codec codecs.Codec, data []byte, object interface{}) error {

	testCodec := new(test.TestCodec)
	service := new(WebCodecService)

	// some test objects
	object := struct{}{}
	data := []byte("Some bytes")

	// setup expectations
	testCodec.On("Unmarshal", data, object).Return(nil)

	// call the target method
	err := service.UnmarshalWithCodec(testCodec, data, object)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, testCodec.Mock)

}

func TestUnmarshalWithCodec_WithError(t *testing.T) {

	// func (s *WebCodecService) UnmarshalWithCodec(codec codecs.Codec, data []byte, object interface{}) error {

	testCodec := new(test.TestCodec)
	service := new(WebCodecService)

	// some test objects
	object := struct{}{}
	data := []byte("Some bytes")

	// setup expectations
	testCodec.On("Unmarshal", data, object).Return(assert.AnError)

	// call the target method
	err := service.UnmarshalWithCodec(testCodec, data, object)

	assert.Equal(t, assert.AnError, err)
	mock.AssertExpectationsForObjects(t, testCodec.Mock)

}