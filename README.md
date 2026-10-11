# go-resume-schema

This repository generates Go types (and optional validation/unmarshal helpers) from the JSON Resume Schema project using `go-jsonschema`.

The generated files are committed to the repo:

- `types.go` (generated from `resume-schema/types.json`)
- `job.go` (generated from `resume-schema/job-schema.json`)
- `resume.go` (generated from `resume-schema/schema.json`)

Hand-written helpers (do not regenerate these):

- `decode.go` — `DecodeResume` / `DecodeJob` accept JSON or YAML. YAML timestamps are rewritten to schema date or RFC3339 strings, then validated through the generated JSON unmarshalers.

## Prerequisites

- Go (see `go.mod`)
- `make`

## Regenerating schema code

```bash
make schema_resume
```

## Pinning the `go-jsonschema` generator version

The Makefile reads the generator version from:

- `go-jsonschema_version.txt`

To bump the generator, edit that file and run:

```bash
make schema_resume
```
