package config

import "path"

const (
	// Artifacts is the artifacts folder
	Artifacts = "/artifacts"
)

var (
	// Report is the report file path
	Report = path.Join(Artifacts, "report.md")
)
