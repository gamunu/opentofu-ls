// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

var (
	version = "0.1.0"
)

// VersionString returns the complete version string, including prerelease
func VersionString() string {
	return version
}
