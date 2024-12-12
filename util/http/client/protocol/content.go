package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type ContentType string

const (
	ContentTypeJSON   ContentType = "application/json"
	HeaderContentType string      = "Content-Type"
)

type PayloadProcessor struct {
	Encoder     BodyEncoder
	Decoder     BodyDecoder
	ContentType ContentType
}

type (
	BodyEncoder func(body interface{}) (*bytes.Buffer, error)
	BodyDecoder func(result interface{}, encoded io.ReadCloser) error
)

func NewPayloadProcessor(typ ContentType) (PayloadProcessor, error) {
	switch typ {
	case ContentTypeJSON:
		return PayloadProcessor{
			Encoder:     jsonEncoder,
			Decoder:     jsonDecoder,
			ContentType: ContentTypeJSON,
		}, nil
	}
	return PayloadProcessor{}, fmt.Errorf("payload type not implemented")
}

func jsonEncoder(data interface{}) (*bytes.Buffer, error) {
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonValue), nil
}

func jsonDecoder(result interface{}, encoded io.ReadCloser) error {
	dec := NewJSONDecoder(encoded)
	return dec.Decode(&result)
}
