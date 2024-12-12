package client

import (
	"fmt"
	"net/http"

	"crossjoin.com/gorxestra/daemon/conductord/api/server/v1/openapi/generated/model"
	"crossjoin.com/gorxestra/data"
	utilClient "crossjoin.com/gorxestra/util/http/client"
)

const (
	healthCheckPath = "health"
	readyCheckPath  = "ready"
	infoCheckPath   = "info"

	registerMusicianPath   = "/v1/musician"
	unregisterMusicianPath = "/v1/musician/%s"
	playMusicPath          = "/v1/music/play/%s"
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

func (h *httpClient) RegisterMusician(m data.Musician) error {
	request := utilClient.Request{
		Path:        registerMusicianPath,
		QueryParams: nil,
		Body: model.Musician{
			Address: m.Address,
			Id:      m.Id.Hex(),
		},
		Method: http.MethodPost,
	}

	return h.restClient.JsonSubmitForm(nil, request)
}

func (h *httpClient) UnregisterMusician(id data.ID) error {
	request := utilClient.Request{
		Path:        fmt.Sprintf(unregisterMusicianPath, id.Hex()),
		QueryParams: nil,
		Body:        nil,
		Method:      http.MethodDelete,
	}

	return h.restClient.JsonSubmitForm(nil, request)
}

func (h *httpClient) PlayMusic(name string) error {
	request := utilClient.Request{
		Path:        fmt.Sprintf(playMusicPath, name),
		QueryParams: nil,
		Body:        nil,
		Method:      http.MethodPost,
	}

	return h.restClient.JsonSubmitForm(nil, request)
}
