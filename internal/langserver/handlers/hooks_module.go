// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-ls/internal/langserver/diagnostics"
	"github.com/hashicorp/terraform-ls/internal/langserver/notifier"
	"github.com/hashicorp/terraform-ls/internal/langserver/session"
	"github.com/hashicorp/terraform-ls/internal/state"
)

func updateDiagnostics(dNotifier *diagnostics.Notifier) notifier.Hook {
	return func(ctx context.Context, changes state.ModuleChanges) error {
		if changes.Diagnostics {
			mod, err := notifier.ModuleFromContext(ctx)
			if err != nil {
				return err
			}

			diags := diagnostics.NewDiagnostics()
			diags.EmptyRootDiagnostic()

			for source, dm := range mod.ModuleDiagnostics {
				diags.Append(source, dm.AutoloadedOnly().AsMap())
			}
			for source, dm := range mod.VarsDiagnostics {
				diags.Append(source, dm.AutoloadedOnly().AsMap())
			}

			dNotifier.PublishHCLDiags(ctx, mod.Path, diags)
		}
		return nil
	}
}

func callRefreshClientCommand(clientRequester session.ClientCaller, commandId string) notifier.Hook {
	return func(ctx context.Context, changes state.ModuleChanges) error {
		// TODO: avoid triggering if module calls/providers did not change
		isOpen, err := notifier.ModuleIsOpen(ctx)
		if err != nil {
			return err
		}

		if isOpen {
			mod, err := notifier.ModuleFromContext(ctx)
			if err != nil {
				return err
			}

			_, err = clientRequester.Callback(ctx, commandId, nil)
			if err != nil {
				return fmt.Errorf("Error calling %s for %s: %s", commandId, mod.Path, err)
			}
		}

		return nil
	}
}

func refreshCodeLens(clientRequester session.ClientCaller) notifier.Hook {
	return func(ctx context.Context, changes state.ModuleChanges) error {
		// TODO: avoid triggering for new targets outside of open module
		if changes.ReferenceOrigins || changes.ReferenceTargets {
			_, err := clientRequester.Callback(ctx, "workspace/codeLens/refresh", nil)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func refreshSemanticTokens(clientRequester session.ClientCaller) notifier.Hook {
	return func(ctx context.Context, changes state.ModuleChanges) error {
		isOpen, err := notifier.ModuleIsOpen(ctx)
		if err != nil {
			return err
		}

		localChanges := isOpen && (changes.TerraformVersion || changes.CoreRequirements ||
			changes.InstalledProviders || changes.ProviderRequirements)

		if localChanges || changes.ReferenceOrigins || changes.ReferenceTargets {
			mod, err := notifier.ModuleFromContext(ctx)
			if err != nil {
				return err
			}

			_, err = clientRequester.Callback(ctx, "workspace/semanticTokens/refresh", nil)
			if err != nil {
				return fmt.Errorf("Error refreshing %s: %s", mod.Path, err)
			}
		}

		return nil
	}
}
