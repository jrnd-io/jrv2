VERSION=0.9.0
GOVERSION=$(shell go version)
USER=$(shell id -u -n)
TIME=$(shell date)
JR_HOME=jr

GOLANCI_LINT_VERSION ?= v1.61.0
GOVULNCHECK_VERSION ?= latest

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
LOCALBIN := $(PROJECT_PATH)/bin

ifndef XDG_DATA_DIRS
ifeq ($(OS), Windows_NT)
	detectedOS := Windows
else
	detectedOS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
endif

ifeq ($(detectedOS), Darwin)
	JR_SYSTEM_DIR="/Library/Application Support"
endif
ifeq ($(detectedOS),  Linux)
	JR_SYSTEM_DIR="/usr/local/share"
endif
ifeq ($(detectedOS), Windows_NT)
	JR_SYSTEM_DIR="$(APPADATA)"
endif
else
	JR_SYSTEM_DIR=$(XDG_DATA_DIRS[0])
endif

ifndef XDG_DATA_HOME
ifeq ($(OS), Windows_NT)
	detectedOS := Windows
else
	detectedOS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
endif

ifeq ($(detectedOS), Darwin)
	JR_USER_DIR="$(HOME)/Library/Application Support"
endif
ifeq ($(detectedOS),  Linux)
	JR_USER_DIR="$(HOME)/.local/share"
endif
ifeq ($(detectedOS), Windows_NT)
	JR_USER_DIR="$(LOCALAPPDATA)"
endif
else
	JR_USER_DIR=$(XDG_DATA_HOME)
endif

hello:
	@echo "JR,the JSON Random Generator"
	@echo " Version: $(VERSION)"
	@echo " Go Version: $(GOVERSION)"
	@echo " Build User: $(USER)"
	@echo " Build Time: $(TIME)"
	@echo " Detected OS: $(detectedOS)"
	@echo " JR System Dir: $(JR_SYSTEM_DIR)"
	@echo " JR User Dir: $(JR_USER_DIR)"

install-gogen:
	go install github.com/actgardner/gogen-avro/v10/cmd/...@latest
	#go install github.com/hamba/avro/v2/cmd/avrogen@latest

generate:
	go generate pkg/generator/generate.go

compile: hello lint test
#compile: hello test
	@echo "Compiling"
	go build -v -ldflags="-s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.GoVersion=$(GOVERSION)' \
	-X 'main.BuildUser=$(USER)' \
	-X 'main.BuildTime=$(TIME)'" \
	-o build/jr github.com/jrnd-io/jrv2/cmd/jr

image:
	docker build -t jrv2:$(VERSION) \
     --build-arg VERSION="$(VERSION)" \
     --build-arg GOVERSION="$(GOVERSION)" \
     --build-arg USER="$(USER)" \
     --build-arg TIME="$(TIME)" \
     .

run: compile
	./build/jr

clean:
	go clean
	rm build/*

test:
	go clean -testcache
	go test --tags=testing ./...

test_coverage:
	go test --tags=testing ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet ./...

.PHONY: lint
lint: golangci-lint
	$(LOCALBIN)/golangci-lint cache clean
	$(LOCALBIN)/golangci-lint run

.PHONY: vuln
vuln: govulncheck
	$(LOCALBIN)/govulncheck -show verbose ./...

.PHONY: check
check: vet lint vuln

help: hello
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}all${RESET}'
	@echo ''

copy_templates:
	sudo mkdir -p $(JR_SYSTEM_DIR)/$(JR_HOME)/kafka && \
	sudo cp -r templates $(JR_SYSTEM_DIR)/$(JR_HOME)
#	cp -r pkg/producers/kafka/*.properties.example $(JR_SYSTEM_DIR)/$(JR_HOME)/kafka/

copy_config:
	mkdir -p $(JR_SYSTEM_DIR)/$(JR_HOME) && \
	cp config/* $(JR_SYSTEM_DIR)/$(JR_HOME)/

install: copy_templates copy_config
	install build/jr /usr/local/bin

all: hello install-gogen generate compile
all_offline: hello generate compile

$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: golangci-lint
golangci-lint: $(LOCALBIN)
	@test -s $(LOCALBIN)/golangci-lint || \
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANCI_LINT_VERSION)

.PHONY: govulncheck
govulncheck: $(LOCALBIN)
	@test -s $(LOCALBIN)/govulncheck || \
	GOBIN=$(LOCALBIN) go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
