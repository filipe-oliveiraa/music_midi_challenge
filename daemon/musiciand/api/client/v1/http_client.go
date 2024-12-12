package client

import (
	"encoding/base64"
	"net/http"

	"crossjoin.com/gorxestra/daemon/musiciand/api/server/v1/openapi/generated/model"
	utilClient "crossjoin.com/gorxestra/util/http/client"
)

const (
	healthCheckPath = "health"
	readyCheckPath  = "ready"
	infoCheckPath   = "info"
	playPath        = "/v1/play"
)

type httpClient struct {
	restClient utilClient.RestClient
}

func (h *httpClient) HealthCheck() (err error) {
	request := utilClient.Request{
		Path:        healthCheckPath,
		QueryParams: nil,
		Body:        nil,
		Method:      http.MethodGet,
	}
	return h.restClient.JsonSubmitForm(nil, request)
}

func (h *httpClient) Ready() (err error) {
	request := utilClient.Request{
		Path:        readyCheckPath,
		QueryParams: nil,
		Body:        nil,
		Method:      http.MethodGet,
	}
	return h.restClient.JsonSubmitForm(nil, request)
}

func (h *httpClient) Info() (resp model.InfoResponse, err error) {
	request := utilClient.Request{
		Path:        infoCheckPath,
		QueryParams: nil,
		Body:        nil,
		Method:      http.MethodGet,
	}
	err = h.restClient.JsonSubmitForm(&resp, request)
	return resp, err
}

func (h *httpClient) Play(bs []byte) error {
	noteBase64 := base64.RawStdEncoding.EncodeToString(bs)

	request := utilClient.Request{
		Path:        playPath,
		QueryParams: nil,
		Body: model.MusicNote{
			Note: noteBase64,
		},
		Method: http.MethodPost,
	}

	return h.restClient.JsonSubmitForm(nil, request)
}
