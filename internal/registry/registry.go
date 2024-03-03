// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

const (
	defaultBaseURL = "https://registry.opentofu.org"
	defaultTimeout = 5 * time.Second
)

type Client struct {
	BaseURL          string
	Timeout          time.Duration
	ProviderPageSize int
	httpClient       *http.Client
}

func NewClient() Client {
	client := cleanhttp.DefaultClient()
	client.Timeout = defaultTimeout

	return Client{
		BaseURL:          defaultBaseURL,
		Timeout:          defaultTimeout,
		ProviderPageSize: 100,
		httpClient:       client,
	}
}
