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

	// Shared schema types are keyed by definition name, not dotted paths.
	for _, typeName := range []string{
		"work", "skills", "program", "location", "link", "dateRange", "exclusion",
	} {
		if _, ok := spec[typeName]; !ok {
			t.Fatalf("resume manifest is missing type %q", typeName)
		}
	}

	// The shared dateRange type lists startDate before endDate.
	dateRange := spec["dateRange"]
	startAt := indexOf(dateRange, "startDate")
	endAt := indexOf(dateRange, "endDate")
	if startAt < 0 || endAt < 0 || startAt >= endAt {
		t.Fatalf("dateRange order = %v, want startDate before endDate", dateRange)
	}

	// Legacy flat pairs on dated sections list startDate before endDate.
	for _, typeName := range []string{"work", "volunteer", "education", "projects"} {
		order := spec[typeName]
		startAt := indexOf(order, "startDate")
		endAt := indexOf(order, "endDate")
		if startAt < 0 || endAt < 0 || startAt >= endAt {
			t.Fatalf("%s date order = %v, want startDate before endDate", typeName, order)
		}
	}
}

// indexOf returns the index of target in list, or -1 when absent.
func indexOf(list []string, target string) int {
	for i, item := range list {
		if item == target {
			return i
		}
	}
	return -1
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

	if _, ok := spec["skills"]; !ok {
		t.Fatal("job manifest is missing type \"skills\"")
	}
}
