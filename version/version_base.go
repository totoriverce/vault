// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package version

import (
	_ "embed"
	"strings"
)

var (
	// The git commit that was compiled. This will be filled in by the compiler.
	GitCommit   string
	GitDescribe string

	// The compilation date. This will be filled in by the compiler.
	BuildDate string

	// Whether cgo is enabled or not; set at build time
	CgoEnabled bool

	// Version info is now being read from the VERSION file and passed in with ldflags,
	// as part of the binary build process in CI.
	// The default values below will be used during local builds.
	//go:embed VERSION
	fullVersion                   string
	Version, VersionPrerelease, _ = strings.Cut(strings.TrimSpace(fullVersion), "-")
	VersionMetadata               = ""
)
