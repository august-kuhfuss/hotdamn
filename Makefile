APP_NAME=hotdamn
BINARY=tmp/$(APP_NAME)
DEBUG=false

.PHONY: default
default: clean install-tools build

.PHONY: clean
clean:
	@rm -rf dist/ tmp/

.PHONY: download
download:
	@echo Download go.mod dependencies
	@go mod download

.PHONY: install-tools
install-tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: update
update:
	@go get -u ./...

.PHONY: build
build:
	@go build -tags='dev' -o $(BINARY) ./cmd/$(APP_NAME)

.PHONY: run
run:
	@\
	bash -c "if [ "$(DEBUG)" = "true" ]; then \
		echo "rebuild"; \
	else \
		$(BINARY) run; \
	fi"

.PHONY: dev
dev:
	@\
	wgo -file .go go mod tidy :: \
	wgo -file .go make -s build :: \
	wgo -file $(BINARY) make -s run

.PHONY: debug
debug:
	@make DEBUG=true dev

.PHONY: snapshot
snapshot:
	@goreleaser release --snapshot --clean
