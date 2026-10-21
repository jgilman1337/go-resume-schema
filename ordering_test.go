package go_resume_schema

import (
	"reflect"
	"testing"
)

// TestResumeOrdering verifies the resume manifest decodes and matches the
// schema.json property order the generator must mirror.
func TestResumeOrdering(t *testing.T) {
	// The manifest parses once and caches.
	spec, err := ResumeOrdering()
	if err != nil {
		t.Fatalf("ResumeOrdering: %v", err)
	}

	// The root order must match schema.json's property order.
	wantRoot := []string{
		"basics", "work", "volunteer", "education", "awards", "certificates",
		"publications", "skills", "languages", "interests", "references",
		"projects", "settings", "meta",
	}
	if got, ok := spec["$"]; !ok || !reflect.DeepEqual(got, wantRoot) {
		t.Fatalf("resume root order = %v, want %v", got, wantRoot)
	}

	// Array items are addressed under their singularized object names.
	for _, path := range []string{"work.work", "skills.skill", "education.education.programs.program"} {
		// A missing entry would leave that object alphabetically sorted.
		if _, ok := spec[path]; !ok {
			t.Fatalf("resume manifest is missing %q", path)
		}
	}
}

// TestJobOrdering verifies the job manifest decodes and matches the
// job-schema.json property order the generator must mirror.
func TestJobOrdering(t *testing.T) {
	// The manifest parses once and caches.
	spec, err := JobOrdering()
	if err != nil {
		t.Fatalf("JobOrdering: %v", err)
	}

	// The root order must match job-schema.json's property order.
	wantRoot := []string{
		"title", "company", "type", "date", "description", "location",
		"remote", "salary", "experience", "responsibilities",
		"qualifications", "optionalQualifications", "keywords", "skills", "meta",
	}
	if got, ok := spec["$"]; !ok || !reflect.DeepEqual(got, wantRoot) {
		t.Fatalf("job root order = %v, want %v", got, wantRoot)
	}

	// The shared location definition is addressed at its document path.
	wantLocation := []string{"address", "postalCode", "city", "countryCode", "region"}
	if got, ok := spec["location"]; !ok || !reflect.DeepEqual(got, wantLocation) {
		t.Fatalf("job location order = %v, want %v", got, wantLocation)
	}
}