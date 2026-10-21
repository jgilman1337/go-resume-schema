// The version file exports the schema release this driver encodes, so
// consumers can stamp documents and compare schema eras without a build
// step or a schema.json parse.
package go_resume_schema

// Version is the resume schema release this driver's types encode. It
// matches this repository's release tag and the schema submodule's
// schema.json version, and moves only on a schema release.
const Version = "2.2.1"
