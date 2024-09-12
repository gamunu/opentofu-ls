// Copyright (c) Gamunu Balagalla.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ProviderList struct {
	Providers []Provider `json:"providers"`
}

type Provider struct {
	Addr               ProviderAddr      `json:"addr"`
	BlockedReason      string            `json:"blocked_reason"`
	CanonicalAddr      ProviderAddr      `json:"canonical_addr"`
	Description        string            `json:"description"`
	ForkCount          int               `json:"fork_count"`
	ForkOf             *ProviderAddr     `json:"fork_of,omitempty"`
	ForkOfLink         string            `json:"fork_of_link"`
	IsBlocked          bool              `json:"is_blocked"`
	Popularity         int               `json:"popularity"`
	ReverseAliases     []ProviderAddr    `json:"reverse_aliases"`
	UpstreamForkCount  int               `json:"upstream_fork_count"`
	UpstreamPopularity int               `json:"upstream_popularity"`
	Versions           []ProviderVersion `json:"versions"`
}

type ProviderAddr struct {
	Display   string `json:"display"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type ProviderVersion struct {
	ID        string    `json:"id"`
	Published time.Time `json:"published"`
}

func (c Client) ListProviders() ([]Provider, error) {
	url := fmt.Sprintf("%s/providers/index.json", c.BaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("unexpected response: %s: %s", resp.Status, string(bodyBytes))
	}

	var response ProviderList
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}

	return response.Providers, nil
}
