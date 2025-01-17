bundle-schemas:
	redocly bundle schemas/odin/openapi.yml -o schemas/odin-schema.yml

lint-schemas:
	redocly lint --config schemas/redocly.yml schemas/odin-schema.yml

split-schemas:
	rm -rf schemas/odin
	redocly split schemas/odin-schema.yml --outDir schemas/odin

generate-api:
	go generate

api-docs:
	redocly preview-docs -p 8048 schemas/odin-schema.yml
