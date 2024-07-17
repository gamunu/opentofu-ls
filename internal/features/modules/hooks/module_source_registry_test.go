// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hooks

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

const responseAWS = `{
	"hits": [
		{
			"full-name": "terraform-aws-modules/vpc/aws",
			"description": "Terraform module which creates VPC resources on AWS",
			"objectID": "modules:23"
		},
		{
			"full-name": "terraform-aws-modules/eks/aws",
			"description": "Terraform module to create an Elastic Kubernetes (EKS) cluster and associated resources",
			"objectID": "modules:1143"
		}
	],
	"nbHits": 10200,
	"page": 0,
	"nbPages": 100,
	"hitsPerPage": 2,
	"exhaustiveNbHits": true,
	"exhaustiveTypo": true,
	"query": "aws",
	"params": "attributesToRetrieve=%5B%22full-name%22%2C%22description%22%5D&hitsPerPage=2&query=aws",
	"renderingContent": {},
	"processingTimeMS": 1,
	"processingTimingsMS": {}
}`

const responseEmpty = `{
	"hits": [],
	"nbHits": 0,
	"page": 0,
	"nbPages": 0,
	"hitsPerPage": 2,
	"exhaustiveNbHits": true,
	"exhaustiveTypo": true,
	"query": "foo",
	"params": "attributesToRetrieve=%5B%22full-name%22%2C%22description%22%5D&hitsPerPage=2&query=foo",
	"renderingContent": {},
	"processingTimeMS": 1
}`

const responseErr = `{
	"message": "Invalid Application-ID or API key",
	"status": 403
}`

type testRequester struct {
	client *http.Client
}

func (r *testRequester) Request(req *http.Request) (*http.Response, error) {
	return r.client.Do(req)
}

// TODO: implement test cases for module search

func buildSearchClientMock(t *testing.T, handler http.HandlerFunc) *search.Client {
	searchServer := httptest.NewTLSServer(handler)
	t.Cleanup(searchServer.Close)

	// Algolia requires hosts to be without a protocol and always assumes https
	u, err := url.Parse(searchServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	searchClient := search.NewClientWithConfig(search.Configuration{
		Hosts: []string{u.Host},
		// We need to disable certificate checking here, because of the self signed cert
		Requester: &testRequester{
			client: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			},
		},
	})

	return searchClient
}
