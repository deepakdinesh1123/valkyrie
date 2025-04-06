sdk-schema-generation-script := "hack/schema/generate_agent_schemas.sh"
sdk-schema-folder := "schemas/agent"

bundle-schemas:
	redocly bundle schemas/valkyrie/openapi.yml -o schemas/valkyrie-schema.yml

lint-schemas:
	redocly lint --config schemas/redocly.yml schemas/valkyrie-schema.yml

split-schemas:
	rm -rf schemas/valkyrie
	redocly split schemas/valkyrie-schema.yml --outDir schemas/valkyrie

api-docs:
	redocly preview-docs -p 8048 schemas/valkyrie-schema.yml
	
generate-json-schema:
	uv run --no-project sdk/valkyriepy/src/generate.py

generate-ts-sdk-schema:
	#!/usr/bin/env sh
	{{sdk-schema-generation-script}} ts {{sdk-schema-folder}} sdk/ts/valkyriejs/schemas

generate-agent-schema:
    #!/usr/bin/env sh
    {{sdk-schema-generation-script}} go {{sdk-schema-folder}} agent/schemas

generate-api-code:
	@just bundle-schemas
	@just lint-schemas
	@go generate

generate-sdk-schemas: 
	@just generate-json-schema
	@just generate-ts-sdk-schema
	@just generate-agent-schema
