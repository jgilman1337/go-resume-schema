package go_resume_schema

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
)

// resumeOrderingJSON is the manifest that declares the output field order of
// every resume schema object. Keys are usually schema type names: "work" for
// work[] items, "dateRange" for dates[] items, "location" for every location
// object, and "$" for the document root. The docorder package resolves these
// keys from document paths and legacy dotted-path entries when needed.
//
//go:embed resume_ordering.json
var resumeOrderingJSON []byte

// jobOrderingJSON is the manifest that declares the output field order of
// every job schema object, using the same dotted-path scheme.
//
//go:embed job_ordering.json
var jobOrderingJSON []byte

// resumeCache holds the resume manifest parsed once for every later caller.
var (
	// resumeOrderingOnce guards the one-time resume manifest parse.
	resumeOrderingOnce sync.Once
	// resumeOrderingSpec holds the parsed resume manifest.
	resumeOrderingSpec map[string][]string
	// resumeOrderingErr holds the parse error, if any.
	resumeOrderingErr error
)

// jobCache holds the job manifest parsed once for every later caller.
var (
	// jobOrderingOnce guards the one-time job manifest parse.
	jobOrderingOnce sync.Once
	// jobOrderingSpec holds the parsed job manifest.
	jobOrderingSpec map[string][]string
	// jobOrderingErr holds the parse error, if any.
	jobOrderingErr error
)

// ResumeOrdering returns the manifest mapping dotted document paths onto the
// field order their objects must carry in backend-generated JSON and YAML.
//
//	@return The manifest mapping paths onto ordered field names.
//	@return An error when the embedded manifest is not valid JSON.
func ResumeOrdering() (map[string][]string, error) {
	// The first caller parses the manifest; every later caller reuses it.
	resumeOrderingOnce.Do(func() {
		resumeOrderingSpec, resumeOrderingErr = parseOrdering(resumeOrderingJSON, "resume_ordering.json")
	})
	return resumeOrderingSpec, resumeOrderingErr
}

// JobOrdering returns the manifest mapping dotted document paths onto the
// field order their objects must carry in backend-generated JSON.
//
//	@return The manifest mapping paths onto ordered field names.
//	@return An error when the embedded manifest is not valid JSON.
func JobOrdering() (map[string][]string, error) {
	// The first caller parses the manifest; every later caller reuses it.
	jobOrderingOnce.Do(func() {
		jobOrderingSpec, jobOrderingErr = parseOrdering(jobOrderingJSON, "job_ordering.json")
	})
	return jobOrderingSpec, jobOrderingErr
}

// parseOrdering decodes one embedded manifest into the spec shape the
// backend's docorder package consumes.
//
//	@param raw The embedded manifest bytes.
//	@param name The manifest file name, for error messages.
//	@return The manifest mapping paths onto ordered field names.
//	@return An error when the embedded manifest is not valid JSON.
func parseOrdering(raw []byte, name string) (map[string][]string, error) {
	// The manifest is a flat map, so one unmarshal captures every entry.
	var spec map[string][]string
	if err := json.Unmarshal(raw, &spec); err != nil {
		return nil, fmt.Errorf("parse %s: %w", name, err)
	}
	return spec, nil
}