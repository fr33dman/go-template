include .env

GOLANGCI_LINT ?= golangci-lint
VERSION := $$(git describe --tags --always)
COMMIT_HASH := $$(git log -1 --pretty=format:"%h")
DATETIME := $$(date -u +%FT%TZ)

.PHONY: init
init: # Init project
	@echo "🛠️ Initializing project..."
	@go mod tidy
	@cp .env.template .env

.PHONY: sqlcgen
sqlcgen: # Generate go queries layer with sqlc from raw sql queries
	@echo "🛠️ Generating sqlc..."
	@find . -name "sqlc.yaml" | while read -r config; do \
		config_dir=$$(dirname "$$config"); \
		root_dir=$$(git rev-parse --show-toplevel); \
		schema_path=$$(python3 -c "import os; print(os.path.relpath(os.path.join('$$root_dir', 'sql', 'migrations'), '$$config_dir'))"); \
		tmp_conf=$$(mktemp "$$config_dir/sqlc.resolved.XXXX.yaml"); \
		yq --yaml-fix-merge-anchor-to-spec=true eval 'explode(.) | with_entries(select(.key | test("^x-") | not))' "$$config" > "$$tmp_conf"; \
		count=$$(yq eval '.sql | length' "$$tmp_conf"); \
		for i in $$(seq 0 $$(($$count - 1))); do \
			existing=$$(yq eval ".sql[$$i].schema // [] | length" "$$tmp_conf"); \
			if [ "$$existing" -eq 0 ]; then \
				echo "Setting schema path in .sql[$$i] to $$schema_path"; \
				yq eval ".sql[$$i].schema = [\"$$schema_path\"]" -i "$$tmp_conf"; \
			fi; \
		done; \
		( cd "$$config_dir" && sqlc generate -f "$$(basename "$$tmp_conf")" ); \
		rm -f "$$tmp_conf"; \
	done

.PHONY: gen
gen: sqlcgen # Generate go code
	@echo "🛠️ Generating go code..."
	@go generate ./gen.go

.PHONY: test
test: # Run tests
	@echo "🧪 Running tests..."
	@go test ./...

.PHONY: coverage
coverage: # Check coverage
	@echo "🧪 Collecting coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

.PHONY: lint
lint: # Run linters
	@echo "🧪 Running linters..."
	@$(GOLANGCI_LINT) run

.PHONY: format
format: # Format go code
	@echo "🧪 Formatting..."
	@gofmt -w $$(find . -name '*.go')

.PHONY: build
build: # Build image
	@echo "🏗️ Building image..."
	@docker build . -f ./build/Dockerfile \
	--build-arg VERSION=$(VERSION) --build-arg COMMIT_HASH=$(COMMIT_HASH) --build-arg DATETIME=$(DATETIME) \
	-t $(IMAGE_NAME):$(VERSION)
