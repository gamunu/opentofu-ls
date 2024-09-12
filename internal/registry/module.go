// Copyright (c) Gamunu Balagalla.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/hashicorp/go-version"
	tfaddr "github.com/opentofu/registry-address"
)

type ClientError struct {
	StatusCode int
	Body       string
}

func (ce ClientError) Error() string {
	return fmt.Sprintf("%d: %s", ce.StatusCode, ce.Body)
}

type Module struct {
	Addr          ModuleAddr                `json:"addr"`
	Description   string                    `json:"description"`
	Versions      []ModuleVersionDescriptor `json:"versions"`
	IsBlocked     bool                      `json:"is_blocked"`
	BlockedReason string                    `json:"blocked_reason"`
}

type ModuleAddr struct {
	Display   string `json:"display"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Target    string `json:"target"`
}

type ModuleVersionDescriptor struct {
	ID        string    `json:"id"`
	Published time.Time `json:"published"`
}

type ModuleVersion struct {
	ID          string    `json:"id"`
	Published   time.Time `json:"published"`
	Description string    `json:"description"`
	Downloads   int       `json:"downloads"`
	Version     string    `json:"version"`
}

type ModuleDetails struct {
	BaseDetails
	Dependencies []ModuleDependency   `json:"dependencies"`
	Providers    []ProviderDependency `json:"providers"`
	Resources    []Resource           `json:"resources"`
}

type BaseDetails struct {
	Readme      bool                `json:"readme"`
	Variables   map[string]Variable `json:"variables"`
	Outputs     map[string]Output   `json:"outputs"`
	SchemaError string              `json:"schema_error"`
	EditLink    string              `json:"edit_link"`
}

type ModuleDependency struct {
	Name              string `json:"name"`
	VersionConstraint string `json:"version_constraint"`
	Source            string `json:"source"`
}

type ProviderDependency struct {
	Alias             string `json:"alias"`
	Name              string `json:"name"`
	FullName          string `json:"full_name"`
	VersionConstraint string `json:"version_constraint"`
}

type Resource struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Name    string `json:"name"`
}

type Variable struct {
	Type        string      `json:"type"`
	Default     interface{} `json:"default"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Sensitive   bool        `json:"sensitive"`
}

type Output struct {
	Description string `json:"description"`
	Sensitive   bool   `json:"sensitive"`
}

func (c Client) GetModuleData(ctx context.Context, addr tfaddr.Module, cons version.Constraints) (*ModuleDetails, error) {
	v, err := c.GetMatchingModuleVersion(ctx, addr, cons)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/modules/%s/%s/%s/%s/index.json",
		c.BaseURL,
		addr.Package.Namespace,
		addr.Package.Name,
		addr.Package.TargetSystem,
		v.String())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, ClientError{StatusCode: resp.StatusCode, Body: string(bodyBytes)}
	}

	var response ModuleDetails
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c Client) GetMatchingModuleVersion(ctx context.Context, addr tfaddr.Module, con version.Constraints) (*version.Version, error) {
	foundVersions, err := c.GetModuleVersions(ctx, addr)
	if err != nil {
		return nil, err
	}

	for _, fv := range foundVersions {
		if con.Check(fv) {
			return fv, nil
		}
	}

	return nil, fmt.Errorf("no suitable version found for %q %q", addr, con)
}

func (c Client) GetModuleVersions(ctx context.Context, addr tfaddr.Module) (version.Collection, error) {
	url := fmt.Sprintf("%s/modules/%s/%s/%s/index.json",
		c.BaseURL,
		addr.Package.Namespace,
		addr.Package.Name,
		addr.Package.TargetSystem)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, ClientError{StatusCode: resp.StatusCode, Body: string(bodyBytes)}
	}

	var module Module
	err = json.NewDecoder(resp.Body).Decode(&module)
	if err != nil {
		return nil, err
	}

	var foundVersions version.Collection
	for _, ver := range module.Versions {
		v, err := version.NewVersion(ver.ID)
		if err == nil {
			foundVersions = append(foundVersions, v)
		}
	}
	sort.Sort(sort.Reverse(foundVersions))

	return foundVersions, nil
}
