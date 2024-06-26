// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/terraform-ls/internal/hooks"
)

func (s *service) AppendCompletionHooks(decoderContext decoder.DecoderContext) {
	h := hooks.Hooks{
		ModStore:       s.modStore,
		RegistryClient: s.registryClient,
		Logger:         s.logger,
	}

	decoderContext.CompletionHooks["CompleteLocalModuleSources"] = h.LocalModuleSources
	decoderContext.CompletionHooks["CompleteRegistryModuleSources"] = h.RegistryModuleSources
	decoderContext.CompletionHooks["CompleteRegistryModuleVersions"] = h.RegistryModuleVersions
}
