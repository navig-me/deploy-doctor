BINARY := deploy-doctor

.PHONY: build test lint coverage test-integration perf-sample

build:
	go build -o bin/$(BINARY) ./cmd/deploy-doctor

test:
	go test ./...

lint:
	go vet ./...

coverage:
	./scripts/check_coverage.sh 60

release-check:
	go test ./...
	go vet ./...
	git diff --quiet && git diff --cached --quiet || (echo "git tree dirty"; exit 1)

test-integration:
	DOCKER_RUNTIME_TEST=1 go test ./test/integration -v

perf-sample: build
	@echo "Run performance sample"
	time ./bin/$(BINARY) scan --static-only
