### SIMPLE MAKEFILE TEMPLATE FOR GO PROJECTS ###

NAME := go-resume-schema
SCHEMA_DIR := ./schema
SCHEMA_AUTOGEN_DIR := ./resume_schema

# ----
# DO NOT EDIT BELOW THIS LINE
# ----
GO := go

# See: https://github.com/omissis/go-jsonschema
GO_JSONSCHEMA := go-jsonschema

.PHONY: schema_resume
schema_resume: REPO_ROOT := https://raw.githubusercontent.com/jgilman1337/resume-schema/refs/heads/master
schema_resume: SCHEMA_GO_PKG := resume_schema
schema_resume: SCHEMA_PKG := resume-schema
schema_resume:
	$(GO_JSONSCHEMA) -e \
		--schema-package=$(REPO_ROOT)/types.json=$(SCHEMA_GO_PKG) \
		--schema-output=$(REPO_ROOT)/types.json=$(SCHEMA_AUTOGEN_DIR)/$(SCHEMA_GO_PKG)/types.go \
		--schema-package=$(REPO_ROOT)/job-schema.json=$(SCHEMA_GO_PKG) \
		--schema-output=$(REPO_ROOT)/job-schema.json=$(SCHEMA_AUTOGEN_DIR)/$(SCHEMA_GO_PKG)/job.go \
		--schema-root-type=$(REPO_ROOT)/job-schema.json=Job \
		--schema-package=$(REPO_ROOT)/schema.json=$(SCHEMA_GO_PKG) \
		--schema-output=$(REPO_ROOT)/schema.json=$(SCHEMA_AUTOGEN_DIR)/$(SCHEMA_GO_PKG)/resume.go \
		--schema-root-type=$(REPO_ROOT)/schema.json=Resume \
		$(SCHEMA_DIR)/$(SCHEMA_PKG)/types.json \
		$(SCHEMA_DIR)/$(SCHEMA_PKG)/job-schema.json \
		$(SCHEMA_DIR)/$(SCHEMA_PKG)/schema.json

.PHONY: check_schema
check_schema:
	./schema/check.sh $(SCHEMA_DIR)/$(SCHEMA_PKG)/schema.json
