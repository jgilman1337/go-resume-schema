package go_resume_schema

import "testing"

// TestEmbeddedSchemasLoad checks that every schema file in SchemaFS decodes.
func TestEmbeddedSchemasLoad(t *testing.T) {
	cases := []struct {
		name string
		load func() (map[string]any, error)
	}{
		{"types", TypesJSONSchema},
		{"job", JobJSONSchema},
		{"resume", ResumeJSONSchema},
		{"keywords", KeywordsDictionaryJSONSchema},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			schema, err := tc.load()
			if err != nil {
				t.Fatal(err)
			}
			if schema["type"] != "object" {
				t.Fatalf("type: %#v", schema["type"])
			}
		})
	}
}

// TestJobJSONSchemaInlinesExperienceEnum checks that experience keeps the closed set from types.json.
func TestJobJSONSchemaInlinesExperienceEnum(t *testing.T) {
	schema, err := JobJSONSchema()
	if err != nil {
		t.Fatalf("JobJSONSchema: %v", err)
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("properties is missing")
	}

	experience, ok := props["experience"].(map[string]any)
	if !ok {
		t.Fatal("experience is missing")
	}
	if _, ok := experience["$ref"]; ok {
		t.Fatal("experience still has a $ref")
	}

	enum, ok := experience["enum"].([]any)
	if !ok {
		t.Fatalf("experience.enum: %#v", experience["enum"])
	}

	found := false
	for _, value := range enum {
		if value == "Junior" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("experience.enum: %#v", enum)
	}
}

// TestResumeJSONSchemaInlinesMeta checks that resume meta is inlined from types.json.
func TestResumeJSONSchemaInlinesMeta(t *testing.T) {
	schema, err := ResumeJSONSchema()
	if err != nil {
		t.Fatalf("ResumeJSONSchema: %v", err)
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("properties is missing")
	}

	meta, ok := props["meta"].(map[string]any)
	if !ok {
		t.Fatal("meta is missing")
	}
	if _, ok := meta["$ref"]; ok {
		t.Fatal("meta still has a $ref")
	}
	if meta["type"] != "object" {
		t.Fatalf("meta.type: %#v", meta["type"])
	}
}
