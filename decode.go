package go_resume_schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"go.yaml.in/yaml/v4"
)

// DecodeResume decodes a Resume from JSON or YAML.
// YAML timestamps are rewritten to schema date or RFC3339 strings before validation.
func DecodeResume(data []byte) (Resume, error) {
	encoded, err := encodeAsJSON(data)
	if err != nil {
		return Resume{}, err
	}

	var resume Resume
	if err := json.Unmarshal(encoded, &resume); err != nil {
		return Resume{}, err
	}

	// An entry mixing both date forms is ambiguous, so it fails decode.
	if err := validateDateForms(&resume); err != nil {
		return Resume{}, err
	}
	return resume, nil
}

// DecodeJob decodes a Job from JSON or YAML.
func DecodeJob(data []byte) (Job, error) {
	encoded, err := encodeAsJSON(data)
	if err != nil {
		return Job{}, err
	}

	var job Job
	if err := json.Unmarshal(encoded, &job); err != nil {
		return Job{}, err
	}
	return job, nil
}

// encodeAsJSON returns JSON bytes for a JSON or YAML document.
func encodeAsJSON(data []byte) ([]byte, error) {
	data = bytes.TrimSpace(data)

	// An empty body is not a schema document.
	if len(data) == 0 {
		return nil, fmt.Errorf("document is empty")
	}

	// JSON already matches the generated UnmarshalJSON validators.
	if isJSON(data) {
		return data, nil
	}

	var raw any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	// A YAML null document has no object to validate.
	if raw == nil {
		return nil, fmt.Errorf("document is empty")
	}

	normalizeYAML(raw)
	encoded, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// isJSON reports whether data looks like a JSON object or array.
func isJSON(data []byte) bool {
	switch data[0] {
	case '{', '[':
		return true
	default:
		return false
	}
}

// normalizeYAML turns YAML timestamps into strings the schema accepts.
func normalizeYAML(value any) {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			// A YAML date node decodes as time.Time, which Iso8601 cannot accept.
			if stamp, ok := child.(time.Time); ok {
				typed[key] = formatSchemaTime(stamp)
				continue
			}
			normalizeYAML(child)
		}
	case []any:
		for i, child := range typed {
			if stamp, ok := child.(time.Time); ok {
				typed[i] = formatSchemaTime(stamp)
				continue
			}
			normalizeYAML(child)
		}
	}
}

// formatSchemaTime writes a date-only string when the clock is midnight, otherwise RFC3339.
func formatSchemaTime(stamp time.Time) string {
	utc := stamp.UTC()

	// Iso8601 fields are year, year-month, or year-month-day.
	if utc.Hour() == 0 && utc.Minute() == 0 && utc.Second() == 0 && utc.Nanosecond() == 0 {
		return utc.Format("2006-01-02")
	}
	return stamp.Format(time.RFC3339)
}

// validateDateForms enforces the dated sections' exclusivity rule: a work,
// volunteer, education, or projects entry carries either the dates array or
// the legacy flat startDate/endDate pair, never both in the same entry. The
// JSON schema cannot spell this exclusion in a form go-jsonschema compiles,
// so decode guards it here.
//
//	@param resume The decoded resume to check.
//	@return An error naming the section and entry that mixed the forms.
func validateDateForms(resume *Resume) error {
	// Work entries may not mix the two date forms.
	for i := range resume.Work {
		// A mixed entry fails decode with its list position.
		if mixesDateForms(resume.Work[i].Dates, resume.Work[i].StartDate, resume.Work[i].EndDate) {
			return fmt.Errorf("work[%d]: an entry carries either dates or the legacy startDate/endDate pair, never both", i)
		}
	}

	// Volunteer entries follow the same rule.
	for i := range resume.Volunteer {
		// A mixed entry fails decode with its list position.
		if mixesDateForms(resume.Volunteer[i].Dates, resume.Volunteer[i].StartDate, resume.Volunteer[i].EndDate) {
			return fmt.Errorf("volunteer[%d]: an entry carries either dates or the legacy startDate/endDate pair, never both", i)
		}
	}

	// Education entries follow the same rule.
	for i := range resume.Education {
		// A mixed entry fails decode with its list position.
		if mixesDateForms(resume.Education[i].Dates, resume.Education[i].StartDate, resume.Education[i].EndDate) {
			return fmt.Errorf("education[%d]: an entry carries either dates or the legacy startDate/endDate pair, never both", i)
		}
	}

	// Projects entries follow the same rule.
	for i := range resume.Projects {
		// A mixed entry fails decode with its list position.
		if mixesDateForms(resume.Projects[i].Dates, resume.Projects[i].StartDate, resume.Projects[i].EndDate) {
			return fmt.Errorf("projects[%d]: an entry carries either dates or the legacy startDate/endDate pair, never both", i)
		}
	}
	return nil
}

// mixesDateForms reports whether one entry carries both date forms: a
// populated dates array alongside any legacy flat startDate or endDate.
//
//	@param dates The entry's dates array.
//	@param start The entry's legacy flat startDate.
//	@param end The entry's legacy flat endDate.
//	@return Whether the two forms appear together.
func mixesDateForms(dates []DateRange, start *Iso8601, end *Iso8601) bool {
	// An absent or empty array cannot clash with the flat pair.
	if len(dates) == 0 {
		return false
	}

	// Any flat key beside a populated array mixes the forms.
	return start != nil || end != nil
}
