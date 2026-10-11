package go_resume_schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
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

// normalizeYAML turns yaml.v3 timestamps into strings the schema accepts.
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
