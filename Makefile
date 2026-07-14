### SIMPLE MAKEFILE TEMPLATE FOR GO PROJECTS ###

NAME := go-resume-schema
SCHEMA_DIR := ./
SCHEMA_AUTOGEN_DIR := ./

# ----
# DO NOT EDIT BELOW THIS LINE
# ----
GO := go

GO_JSONSCHEMA_VERSION_FILE := go-jsonschema_version.txt
# Keep the generator version configurable without editing this Makefile.
# If the file doesn't exist, fall back to the previous default.
GO_JSONSCHEMA_VERSION := $(strip $(shell \
	if [ -f "$(GO_JSONSCHEMA_VERSION_FILE)" ]; then \
		tr -d '\n' < "$(GO_JSONSCHEMA_VERSION_FILE)"; \
	else \
		echo v0.23.1; \
	fi \
))

# This Makefile automatically pulls in the required version of go-jsonschema for use in this project.
# See: https://github.com/omissis/go-jsonschema
GO_JSONSCHEMA := go run github.com/atombender/go-jsonschema@$(GO_JSONSCHEMA_VERSION)

.PHONY: schema_resume
schema_resume: REPO_ROOT := https://raw.githubusercontent.com/jgilman1337/resume-schema/refs/heads/master
schema_resume: SCHEMA_GO_PKG := go_resume_schema
schema_resume: SCHEMA_PKG := resume-schema
schema_resume:
	$(GO_JSONSCHEMA) -e \
		--schema-package=$(REPO_ROOT)/types.json=$(SCHEMA_GO_PKG) \
		--schema-output=$(REPO_ROOT)/types.json=$(SCHEMA_AUTOGEN_DIR)/types.go \
		--schema-package=$(REPO_ROOT)/job-schema.json=$(SCHEMA_GO_PKG) \
		--schema-output=$(REPO_ROOT)/job-schema.json=$(SCHEMA_AUTOGEN_DIR)/job.go \
		--schema-root-type=$(REPO_ROOT)/job-schema.json=Job \
		--schema-package=$(REPO_ROOT)/schema.json=$(SCHEMA_GO_PKG) \
		--schema-output=$(REPO_ROOT)/schema.json=$(SCHEMA_AUTOGEN_DIR)/resume.go \
		--schema-root-type=$(REPO_ROOT)/schema.json=Resume \
		$(SCHEMA_DIR)/$(SCHEMA_PKG)/types.json \
		$(SCHEMA_DIR)/$(SCHEMA_PKG)/job-schema.json \
		$(SCHEMA_DIR)/$(SCHEMA_PKG)/schema.json
