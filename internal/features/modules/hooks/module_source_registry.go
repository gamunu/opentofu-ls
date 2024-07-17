// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hooks

import (
	"context"
	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/zclconf/go-cty/cty"
	"strings"
)

type RegistryModule struct {
	FullName    string `json:"full-name"`
	Description string `json:"description"`
}

func (h *Hooks) RegistryModuleSources(ctx context.Context, value cty.Value) ([]decoder.Candidate, error) {
	candidates := make([]decoder.Candidate, 0)
	prefix := value.AsString()

	if strings.HasPrefix(prefix, ".") {
		// We're likely dealing with a local module source here; no need to search the registry
		// A search for "." will not return any results
		return candidates, nil
	}

	//TODO: figure out how we can do module autocompletion

	return candidates, nil
}
