// Copyright (c) Gamunu Balagalla.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	tfaddr "github.com/opentofu/registry-address"
)

func TestGetModuleData(t *testing.T) {
	ctx := context.Background()
	addr, err := tfaddr.ParseModuleSource("azure/alz/azurerm")
	if err != nil {
		t.Fatal(err)
	}

	cons := version.MustConstraints(version.NewConstraint("0.8.1"))

	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/modules/azure/alz/azurerm/index.json":
			w.Write([]byte(moduleDataMockResponse))
		case "/modules/azure/alz/azurerm/v0.8.1/index.json":
			w.Write([]byte(moduleVersionsMockResponse))
		default:
			http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
		}
	}))
	client.BaseURL = srv.URL
	t.Cleanup(srv.Close)

	data, err := client.GetModuleData(ctx, addr, cons)
	if err != nil {
		t.Fatal(err)
	}

	expectedData := &ModuleDetails{
		ID:        "v0.8.1",
		Published: time.Date(2024, time.August, 5, 17, 16, 42, 0, time.FixedZone("", 3600)),
		Readme:    true,
		Inputs: map[string]Variable{
			"architecture_name": {
				Type:        "string",
				Description: "The name of the architecture to create. This needs to be*.alz_architecture_definition.[json|yaml|yml] files.\n",
				Required:    true,
			},
			"location": {
				Type:        "string",
				Description: "The default location for resources in this management group. Used for policy managed identities.\n",
				Required:    true,
			},
		},
		Outputs: map[string]Output{
			"management_group_resource_ids": {
				Description: "A map of management group names to their resource ids.",
				Sensitive:   false,
			},
			"policy_assignment_resource_ids": {
				Description: "A map of policy assignment names to their resource ids.",
				Sensitive:   false,
			},
		},
		Providers: []ProviderDependency{},
		Dependencies: []ModuleDependency{
			{
				Name:              "policy_assignment",
				VersionConstraint: "",
				Source:            "./modules/azapi_helper",
			},
			{
				Name:              "policy_definitions",
				VersionConstraint: "",
				Source:            "./modules/azapi_helper",
			},
		},
		Resources: []Resource{
			{
				Address: "modtm_telemetry.telemetry",
				Type:    "modtm_telemetry",
				Name:    "telemetry",
			},
			{
				Address: "random_uuid.telemetry",
				Type:    "random_uuid",
				Name:    "telemetry",
			},
		},
		Submodules: map[string]Submodule{
			"azapi_helper": {
				ModuleDetails: ModuleDetails{
					Readme: true,
					Inputs: map[string]Variable{
						"body": {
							Type:        "dynamic",
							Description: "The body object of the resource.",
							Required:    true,
						},
						"name": {
							Type:        "string",
							Description: "The name of resource.",
							Required:    true,
						},
					},
					Outputs: map[string]Output{
						"identity": {
							Description: "The identity configuration of the resource.",
							Sensitive:   false,
						},
						"name": {
							Description: "The name of the resource.",
							Sensitive:   false,
						},
					},
					Providers:    []ProviderDependency{},
					Dependencies: []ModuleDependency{},
					Resources: []Resource{
						{
							Address: "azapi_resource.this",
							Type:    "azapi_resource",
							Name:    "this",
						},
						{
							Address: "terraform_data.replace_trigger",
							Type:    "terraform_data",
							Name:    "replace_trigger",
						},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(expectedData, data); diff != "" {
		t.Fatalf("mismatched data: %s", diff)
	}
}

func TestGetMatchingModuleVersion(t *testing.T) {
	ctx := context.Background()
	addr, err := tfaddr.ParseModuleSource("azure/alz/azurerm")
	if err != nil {
		t.Fatal(err)
	}
	cons := version.MustConstraints(version.NewConstraint(">=0.7.0"))
	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/modules/azure/alz/azurerm/index.json" {
			w.Write([]byte(moduleDataMockResponse))
			return
		}
		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	client.BaseURL = srv.URL
	t.Cleanup(srv.Close)

	v, err := client.GetMatchingModuleVersion(ctx, addr, cons)
	if err != nil {
		t.Fatal(err)
	}

	expectedVersion := version.Must(version.NewVersion("v0.8.1"))
	if !expectedVersion.Equal(v) {
		t.Fatalf("expected version: %s, given: %s", expectedVersion, v)
	}
}

func TestGetModuleVersions(t *testing.T) {
	ctx := context.Background()
	addr, err := tfaddr.ParseModuleSource("azure/alz/azurerm")
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/modules/azure/alz/azurerm/index.json" {
			w.Write([]byte(moduleDataMockResponse))
			return
		}
		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	client.BaseURL = srv.URL
	t.Cleanup(srv.Close)

	versions, err := client.GetModuleVersions(ctx, addr)
	if err != nil {
		t.Fatal(err)
	}

	expectedVersions := version.Collection{
		version.Must(version.NewVersion("0.8.1")),
		version.Must(version.NewVersion("0.8.0")),
		version.Must(version.NewVersion("0.7.0")),
		version.Must(version.NewVersion("0.6.0")),
		version.Must(version.NewVersion("0.5.0")),
		version.Must(version.NewVersion("0.4.1")),
		version.Must(version.NewVersion("0.4.0")),
		version.Must(version.NewVersion("0.3.3")),
		version.Must(version.NewVersion("0.3.2")),
		version.Must(version.NewVersion("0.3.1")),
		version.Must(version.NewVersion("0.3.0")),
		version.Must(version.NewVersion("0.2.0")),
		version.Must(version.NewVersion("0.1.1")),
		version.Must(version.NewVersion("0.1.0")),
	}

	if diff := cmp.Diff(expectedVersions, versions); diff != "" {
		t.Fatalf("mismatched versions: %s", diff)
	}
}

func TestCancellationThroughContext(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, 50*time.Millisecond)
	t.Cleanup(cancelFunc)

	addr, err := tfaddr.ParseModuleSource("azure/alz/azurerm")
	if err != nil {
		t.Fatal(err)
	}
	cons := version.MustConstraints(version.NewConstraint(">=0.7.0"))
	client := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond) // Delay longer than the context timeout
		if r.RequestURI == "/modules/azure/alz/azurerm/index.json" {
			w.Write([]byte(moduleDataMockResponse))
			return
		}
		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	client.BaseURL = srv.URL
	t.Cleanup(srv.Close)

	_, err = client.GetMatchingModuleVersion(ctx, addr, cons)
	if err == nil {
		t.Fatal("expected error due to context cancellation, got nil")
	}

	urlErr, ok := err.(*url.Error)
	if !ok {
		t.Fatalf("expected *url.Error, got: %T", err)
	}

	if urlErr.Err != context.DeadlineExceeded {
		t.Fatalf("expected context.DeadlineExceeded error, got: %v", urlErr.Err)
	}
}

func TestClientError(t *testing.T) {
	err := ClientError{
		StatusCode: 404,
		Body:       "Not Found",
	}

	expected := "404: Not Found"
	if err.Error() != expected {
		t.Fatalf("expected error message: %s, got: %s", expected, err.Error())
	}
}
