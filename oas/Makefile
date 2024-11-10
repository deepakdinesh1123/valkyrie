
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