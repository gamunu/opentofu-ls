// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hooks

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	tfmod "github.com/gamunu/opentofu-schema/module"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-ls/internal/features/modules/state"
	"github.com/hashicorp/terraform-ls/internal/registry"
	globalState "github.com/hashicorp/terraform-ls/internal/state"
	tfaddr "github.com/opentofu/registry-address"
	"github.com/zclconf/go-cty/cty"
)

var moduleVersionsMockResponse = `{
  "addr": {
    "display": "azure/aks/azurerm",
    "namespace": "azure",
    "name": "aks",
    "target": "azurerm"
  },
  "description": "Terraform Module for deploying an AKS cluster",
  "versions": [
    {
      "id": "v9.1.0",
      "published": "2024-07-04T07:12:29+01:00"
    },
    {
      "id": "v9.0.0",
      "published": "2024-06-07T02:31:28+01:00"
    },
    {
      "id": "v8.0.0",
      "published": "2024-03-05T07:33:07Z"
    }
  ],
  "is_blocked": false
}`

func TestHooks_RegistryModuleVersions(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	ctx = decoder.WithPath(ctx, lang.Path{
		Path:       tmpDir,
		LanguageID: "terraform",
	})
	ctx = decoder.WithPos(ctx, hcl.Pos{
		Line:   2,
		Column: 5,
		Byte:   5,
	})
	ctx = decoder.WithFilename(ctx, "main.tf")
	ctx = decoder.WithMaxCandidates(ctx, 3)
	s, err := globalState.NewStateStore()
	if err != nil {
		t.Fatal(err)
	}
	store, err := state.NewModuleStore(s.ProviderSchemas, s.RegistryModules, s.ChangeStore)
	if err != nil {
		t.Fatal(err)
	}

	regClient := registry.NewClient()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/modules/azure/aks/azurerm/index.json" {
			w.Write([]byte(moduleVersionsMockResponse))
			return
		}
		http.Error(w, fmt.Sprintf("unexpected request: %q", r.RequestURI), 400)
	}))
	regClient.BaseURL = srv.URL
	t.Cleanup(srv.Close)

	h := &Hooks{
		ModStore:       store,
		RegistryClient: regClient,
	}

	err = store.Add(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	metadata := &tfmod.Meta{
		Path: tmpDir,
		ModuleCalls: map[string]tfmod.DeclaredModuleCall{
			"vpc": {
				LocalName:  "vpc",
				SourceAddr: tfaddr.MustParseModuleSource("registry.opentofu.org/azure/aks/azurerm"),
				RangePtr: &hcl.Range{
					Filename: "main.tf",
					Start:    hcl.Pos{Line: 1, Column: 1, Byte: 1},
					End:      hcl.Pos{Line: 4, Column: 2, Byte: 20},
				},
			},
		},
	}
	err = store.UpdateMetadata(tmpDir, metadata, nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedCandidates := []decoder.Candidate{
		{
			Label:         `"9.1.0"`,
			Kind:          lang.StringCandidateKind,
			RawInsertText: `"9.1.0"`,
			SortText:      "  0",
		},
		{
			Label:         `"9.0.0"`,
			Kind:          lang.StringCandidateKind,
			RawInsertText: `"9.0.0"`,
			SortText:      "  1",
		},
		{
			Label:         `"8.0.0"`,
			Kind:          lang.StringCandidateKind,
			RawInsertText: `"8.0.0"`,
			SortText:      "  2",
		},
	}

	candidates, _ := h.RegistryModuleVersions(ctx, cty.StringVal(""))
	fmt.Print(candidates)
	if diff := cmp.Diff(expectedCandidates, candidates); diff != "" {
		t.Fatalf("mismatched candidates: %s", diff)
	}
}
