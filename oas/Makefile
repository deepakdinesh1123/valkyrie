
.PHONY: bundle-api-spec
bundle-api-spec:
	redocly bundle oas/odin/openapi.yml -o oas/odin-schema.yml

.PHONY: oaslint
oaslint:
	redocly lint --config oas/redocly.yml oas/odin-schema.yml