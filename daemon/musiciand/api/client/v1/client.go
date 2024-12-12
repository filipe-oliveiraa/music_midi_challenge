package client

import (
	"fmt"
	"net/url"

	"crossjoin.com/gorxestra/daemon/conductord/api"
	utilClient "crossjoin.com/gorxestra/util/http/client"
)

type ClientDaemon interface {
	api.NodeInterface
}

type client struct {
	*httpClient
}

func New(domain string) (client, error) {
	u, err := url.Parse(domain)
	if err != nil {
		return client{}, fmt.Errorf("parsing url: %w", err)
	}

	c := client{
		httpClient: &httpClient{
			restClient: utilClient.MakeRestClient(*u),
		},
	}
	return c, nil
}
