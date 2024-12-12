package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"crossjoin.com/gorxestra/util/http/client/protocol"
	"crossjoin.com/gorxestra/util/http/common"
	"crossjoin.com/gorxestra/util/http/query"
)

const (
	maxRawResponseBytes = 50e6
)

// unauthorizedRequestError is generated when we receive 401 error from the server. This error includes the inner error
// as well as the likely parameters that caused the issue.
type unauthorizedRequestError struct {
	errorString string
	url         string
}

// Error format an error string for the unauthorizedRequestError error.
func (e unauthorizedRequestError) Error() string {
	return fmt.Sprintf("Unauthorized request to `%s` when using : %s", e.url, e.errorString)
}

// HTTPError is generated when we receive an unhandled error from the server. This error contains the error string.
type HTTPError struct {
	StatusCode  int
	Status      string
	ErrorString string
	Data        map[string]any
}

// Error formats an error string.
func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP %s: %s", e.Status, e.ErrorString)
}

type ErrorMapper func(errStr string) (error, bool)

// RestClient manages the REST interface for a calling user.
type RestClient struct {
	serverURL url.URL
	errMapper ErrorMapper
}

// MakeRestClient is the factory for constructing a RestClient for a given endpoint
func MakeRestClient(url url.URL) RestClient {
	return MakeRestClientWithMapper(url, func(errStr string) (error, bool) {
		return nil, false
	})
}

// MakeRestClientWithMapper is the factory for constructing a RestClient for a given endpoint
// receiving a error mapper to transform string errors into errors
func MakeRestClientWithMapper(url url.URL, mapper ErrorMapper) RestClient {
	return RestClient{
		serverURL: url,
		errMapper: mapper,
	}
}

// filterASCII filter out the non-ascii printable characters out of the given input string.
// It's used as a security qualifier before adding network provided data into an error message.
// The function allows only characters in the range of [32..126], which excludes all the
// control character, new lines, deletion, etc. All the alpha numeric and punctuation characters
// are included in this range.
func filterASCII(unfilteredString string) (filteredString string) {
	for i, r := range unfilteredString {
		if int(r) >= 0x20 && int(r) <= 0x7e {
			filteredString += string(unfilteredString[i])
		}
	}
	return
}

// extractError checks if the response signifies an error (for now, StatusCode != 200 or StatusCode != 201).
// If so, it returns the error.
// Otherwise, it returns nil.
func extractError(mapper ErrorMapper, resp *http.Response) error {
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	}

	errorBuf, _ := io.ReadAll(resp.Body) // ignore returned error
	var errorJSON common.Error
	decodeErr := json.Unmarshal(errorBuf, &errorJSON)

	var errorString string
	var data map[string]any
	if decodeErr == nil {
		errorString = errorJSON.Error
	} else {
		errorString = string(errorBuf)
	}
	errorString = filterASCII(errorString)

	if resp.StatusCode == http.StatusUnauthorized {
		return unauthorizedRequestError{errorString, resp.Request.URL.String()}
	}

	err, ok := mapper(errorString)
	if ok {
		return err
	}

	return HTTPError{
		StatusCode:  resp.StatusCode,
		Status:      resp.Status,
		ErrorString: errorString,
		Data:        data,
	}
}

// RawResponse is fulfilled by responses that should not be decoded as json
type RawResponse interface {
	SetBytes([]byte)
}

// mergeRawQueries merges two raw queries, appending an "&" if both are non-empty
func mergeRawQueries(q1, q2 string) string {
	if q1 == "" || q2 == "" {
		return q1 + q2
	}
	return q1 + "&" + q2
}

type Request struct {
	Path        string
	QueryParams interface{}
	Body        interface{}
	Method      string
}

// submitForm is a helper used for submitting (ex.) GETs and POSTs to the server
// if expectNoContent is true, then it is expected that the response received will have a content length of zero
//
//nolint:funlen
func (client RestClient) submitForm(
	response interface{}, request Request,
	payloadProcessor protocol.PayloadProcessor, expectNoContent bool,
) error {
	var err error
	queryURL := client.serverURL
	queryURL.Path, err = url.JoinPath(queryURL.Path, request.Path)
	if err != nil {
		return err
	}

	var req *http.Request
	var v url.Values

	if request.QueryParams != nil {
		v, err = query.Values(request.QueryParams)
		if err != nil {
			return err
		}
	}

	bodyReader, err := payloadProcessor.Encoder(request.Body)
	if err != nil {
		return err
	}

	queryURL.RawQuery = mergeRawQueries(queryURL.RawQuery, v.Encode())
	req, err = http.NewRequest(request.Method, queryURL.String(), bodyReader)
	if err != nil {
		return err
	}

	if payloadProcessor.ContentType != "" {
		req.Header.Add(protocol.HeaderContentType, string(payloadProcessor.ContentType))
	}

	httpClient := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	defer httpClient.CloseIdleConnections()

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	// Ensure response isn't too large
	resp.Body = http.MaxBytesReader(nil, resp.Body, maxRawResponseBytes)
	defer resp.Body.Close()

	err = extractError(client.errMapper, resp)
	if err != nil {
		return err
	}

	if expectNoContent {
		if resp.ContentLength == 0 {
			return nil
		}
		return fmt.Errorf("expected empty response but got response of %d bytes", resp.ContentLength)
	}

	if response != nil {
		err = payloadProcessor.Decoder(&response, resp.Body)
	}

	if err != nil {
		return err
	}
	return nil
}

// get performs a request in format json against the path
func (client RestClient) JsonSubmitForm(response interface{}, request Request) error {
	p, err := protocol.NewPayloadProcessor(protocol.ContentTypeJSON)
	if err != nil {
		return err
	}
	return client.submitForm(response, request, p, false)
}
