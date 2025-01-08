SDK_SCHEMA_GENERATION_SCRIPT = hack/schema/generate_agent_schemas.sh
SDK_SCHEMA_FOLDER = schemas/agent

.PHONY: bundle-oas
bundle-oas:
	redocly bundle oas/odin/openapi.yml -o oas/odin-schema.yml

.PHONY: lint-oas
lint-oas:
	redocly lint --config oas/redocly.yml oas/odin-schema.yml

.PHONY: split-oas
split-oas:
	rm -rf oas/odin
	redocly split oas/odin-schema.yml --outDir oas/odin

.PHONY: generate-api
generate-api:
	go generate

.PHONY: api-docs
api-docs:
	redocly preview-docs -p 8048 oas/odin-schema.yml

.PHONY: generate-json-schema
generate-json-schema:
	uv run --no-project sdk/python/pyvalkyrie/schema.py

.PHONY: generate-js-sdk-schema
generate-js-sdk-schema:
	quicktype -o sdk/js/valkyriejs/schema.ts --src-lang schema --lang typescript-zod schemas/agent/execute_command.json

.PHONY: generate-agent-schema
generate-agent-schema:
	@$(SDK_SCHEMA_GENERATION_SCRIPT) go $(SDK_SCHEMA_FOLDER) agent/schema

.PHONY: generate-sdk-schemas
generate-sdk-schemas: generate-json-schema generate-agent-schema generate-js-sdk-schema 
	@echo "Generating SDK Schemas"


