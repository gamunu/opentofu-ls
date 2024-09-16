// Copyright (c) Gamunu Balagalla.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestListProviders(t *testing.T) {
	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/providers/index.json" {
			w.Write([]byte(`{
  "providers": [
    {
      "addr": {
        "display": "hashicorp/azuread",
        "namespace": "hashicorp",
        "name": "azuread"
      },
      "description": "Manage Azure Active Directory resources",
      "popularity": 1000,
      "fork_count": 50,
      "versions": [
        {
          "id": "0.36.1",
          "published": "2022-08-24T19:09:29Z"
        }
      ],
      "is_blocked": false
    },
    {
      "addr": {
        "display": "hashicorp/http",
        "namespace": "hashicorp",
        "name": "http"
      },
      "description": "Interact with HTTP servers",
      "popularity": 500,
      "fork_count": 20,
      "versions": [
        {
          "id": "3.2.1",
          "published": "2022-06-21T18:56:56Z"
        }
      ],
      "is_blocked": false
    }
  ]
}`))
			return
		}
		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	client.BaseURL = srv.URL
	t.Cleanup(srv.Close)

	providers, err := client.ListProviders()
	if err != nil {
		t.Fatal(err)
	}

	expectedProviders := []Provider{
		{
			Addr: ProviderAddr{
				Display:   "hashicorp/azuread",
				Name:      "azuread",
				Namespace: "hashicorp",
			},
			Description: "Manage Azure Active Directory resources",
			Popularity:  1000,
			ForkCount:   50,
			Versions: []ProviderVersion{
				{
					ID:        "0.36.1",
					Published: time.Date(2022, 8, 24, 19, 9, 29, 0, time.UTC),
				},
			},
			IsBlocked: false,
		},
		{
			Addr: ProviderAddr{
				Display:   "hashicorp/http",
				Name:      "http",
				Namespace: "hashicorp",
			},
			Description: "Interact with HTTP servers",
			Popularity:  500,
			ForkCount:   20,
			Versions: []ProviderVersion{
				{
					ID:        "3.2.1",
					Published: time.Date(2022, 6, 21, 18, 56, 56, 0, time.UTC),
				},
			},
			IsBlocked: false,
		},
	}
	if diff := cmp.Diff(expectedProviders, providers); diff != "" {
		t.Fatalf("unexpected providers: %s", diff)
	}
}
