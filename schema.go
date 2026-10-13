package go_resume_schema

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

// SchemaFS holds the JSON Schema documents this package is generated from.
//
//go:embed resume-schema/types.json resume-schema/job-schema.json resume-schema/schema.json resume-schema/keywords-dictionary-schema.json
var SchemaFS embed.FS

const (
	// TypesSchemaPath is types.json inside SchemaFS.
	TypesSchemaPath = "resume-schema/types.json"
	// JobSchemaPath is job-schema.json inside SchemaFS.
	JobSchemaPath = "resume-schema/job-schema.json"
	// ResumeSchemaPath is schema.json inside SchemaFS.
	ResumeSchemaPath = "resume-schema/schema.json"
	// KeywordsDictionarySchemaPath is keywords-dictionary-schema.json inside SchemaFS.
	KeywordsDictionarySchemaPath = "resume-schema/keywords-dictionary-schema.json"

	typesRefPrefix = "types.json#/definitions/"
	localRefPrefix = "#/definitions/"
)

// TypesJSONSchema returns the shared type definitions document.
func TypesJSONSchema() (map[string]any, error) {
	return readJSONObject(TypesSchemaPath)
}

// JobJSONSchema returns job-schema.json with types.json $ref values inlined.
func JobJSONSchema() (map[string]any, error) {
	return resolvedWithTypes(JobSchemaPath)
}

// ResumeJSONSchema returns schema.json with types.json $ref values inlined.
func ResumeJSONSchema() (map[string]any, error) {
	return resolvedWithTypes(ResumeSchemaPath)
}

// KeywordsDictionaryJSONSchema returns the keywords dictionary schema document.
func KeywordsDictionaryJSONSchema() (map[string]any, error) {
	return readJSONObject(KeywordsDictionarySchemaPath)
}

// resolvedWithTypes loads a schema file and inlines types.json definitions.
func resolvedWithTypes(path string) (map[string]any, error) {
	doc, err := readJSONObject(path)
	if err != nil {
		return nil, err
	}

	defs, err := typesDefinitions()
	if err != nil {
		return nil, err
	}

	resolved, err := resolveRefs(doc, defs)
	if err != nil {
		return nil, err
	}

	schema, ok := resolved.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s is not an object", path)
	}

	return schema, nil
}

// typesDefinitions returns the definitions object from types.json.
func typesDefinitions() (map[string]any, error) {
	typesDoc, err := TypesJSONSchema()
	if err != nil {
		return nil, err
	}

	defs, _ := typesDoc["definitions"].(map[string]any)
	if defs == nil {
		return nil, fmt.Errorf("types schema is missing definitions")
	}

	return defs, nil
}

// readJSONObject loads an embedded JSON object.
func readJSONObject(name string) (map[string]any, error) {
	raw, err := SchemaFS.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", name, err)
	}

	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("decode %s: %w", name, err)
	}
	if doc == nil {
		return nil, fmt.Errorf("%s is empty", name)
	}

	return doc, nil
}

// resolveRefs replaces JSON Schema $ref nodes with the referenced definition.
func resolveRefs(node any, defs map[string]any) (any, error) {
	switch typed := node.(type) {
	case map[string]any:
		// A $ref node is replaced by the definition it names.
		ref, ok := typed["$ref"].(string)
		if ok {
			return resolveRefNode(typed, ref, defs)
		}

		out := make(map[string]any, len(typed))
		for key, child := range typed {
			resolved, err := resolveRefs(child, defs)
			if err != nil {
				return nil, err
			}
			out[key] = resolved
		}
		return out, nil

	case []any:
		out := make([]any, len(typed))
		for i, child := range typed {
			resolved, err := resolveRefs(child, defs)
			if err != nil {
				return nil, err
			}
			out[i] = resolved
		}
		return out, nil

	default:
		return node, nil
	}
}

// resolveRefNode inlines a $ref and overlays sibling keywords such as description.
func resolveRefNode(node map[string]any, ref string, defs map[string]any) (any, error) {
	name, err := definitionName(ref)
	if err != nil {
		return nil, err
	}

	def, ok := defs[name]
	if !ok {
		return nil, fmt.Errorf("unknown schema ref %q", ref)
	}

	resolved, err := resolveRefs(def, defs)
	if err != nil {
		return nil, err
	}

	base, ok := resolved.(map[string]any)
	// A non-object definition has no sibling keywords to merge.
	if !ok {
		return resolved, nil
	}

	out := make(map[string]any, len(base)+len(node)-1)
	for key, value := range base {
		out[key] = value
	}
	for key, value := range node {
		// $ref is the pointer that was just inlined.
		if key == "$ref" {
			continue
		}
		out[key] = value
	}

	return out, nil
}

// definitionName extracts the definition key from a types.json or local $ref.
func definitionName(ref string) (string, error) {
	switch {
	case strings.HasPrefix(ref, typesRefPrefix):
		return strings.TrimPrefix(ref, typesRefPrefix), nil
	case strings.HasPrefix(ref, localRefPrefix):
		return strings.TrimPrefix(ref, localRefPrefix), nil
	default:
		return "", fmt.Errorf("unsupported schema ref %q", ref)
	}
}
