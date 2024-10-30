include testing/benchmark/execution/loadtest.env
export

.PHONY: generate-load-test

generate-load-test:
	@openapi-generator-cli generate \
		-i "oas/odin-schema.yml" \
		-g k6 \
		-o "testing/benchmark/execution"

.PHONY: generate-load-test-data

generate-load-test-data:
	go run testing/benchmark/execution/data.go

.PHONY: run-load-test

run-load-test:
	./k6 run testing/benchmark/execution/script.js

get-k6-binary:
	docker run --rm -it -u "$(id -u):$(id -g)" -v "${PWD}:/xk6" grafana/xk6 build v0.45.1 \
		--with github.com/phymbert/xk6-sse@v0.1.2