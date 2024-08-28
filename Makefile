VERSION=2.0.0
GOVERSION=$(shell go version)
USER=$(shell id -u -n)
TIME=$(shell date)
JR_HOME=jr

ifndef XDG_DATA_DIRS
ifeq ($(OS), Windows_NT)
	detectedOS := Windows
else
	detectedOS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
endif

ifeq ($(detectedOS), Darwin)
	JR_SYSTEM_DIR="$(HOME)/Library/Application Support"
endif
ifeq ($(detectedOS),  Linux)
	JR_SYSTEM_DIR="$(HOME)/.config"
endif
ifeq ($(detectedOS), Windows_NT)
	JR_SYSTEM_DIR="$(LOCALAPPDATA)"
endif
else
	JR_SYSTEM_DIR=$(XDG_DATA_DIRS)
endif

ifndef XDG_DATA_HOME
ifeq ($(OS), Windows_NT)
	detectedOS := Windows
else
	detectedOS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
endif

ifeq ($(detectedOS), Darwin)
	JR_USER_DIR="$(HOME)/.local/share"
endif
ifeq ($(detectedOS),  Linux)
	JR_USER_DIR="$(HOME)/.local/share"
endif
ifeq ($(detectedOS), Windows_NT)
	JR_USER_DIR="$(LOCALAPPDATA)" //@TODO
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
	@echo "Compiling"
	go build -v -ldflags="-X 'main.Version=$(VERSION)' \
	-X 'main.GoVersion=$(GOVERSION)' \
	-X 'main.BuildUser=$(USER)' \
	-X 'main.BuildTime=$(TIME)'" \
	-o build/jr github.com/jrnd-io/jrv2/cmd/jr

run: compile
	./build/jr

clean:
	go clean
	rm build/*

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --config .localci/lint/golangci.yml

help: hello
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}all${RESET}'
	@echo ''

copy_templates:
	mkdir -p $(JR_SYSTEM_DIR)/$(JR_HOME)/kafka && \
	cp -r templates $(JR_SYSTEM_DIR)/$(JR_HOME) && \
	cp -r pkg/producers/kafka/*.properties.example $(JR_SYSTEM_DIR)/$(JR_HOME)/kafka/

copy_config:
	mkdir -p $(JR_SYSTEM_DIR)/$(JR_HOME) && \
	cp config/* $(JR_SYSTEM_DIR)/$(JR_HOME)/

install:
	install build/jr /usr/local/bin

test:
	go clean -testcache
	go test ./...

all: hello install-gogen generate compile
all_offline: hello generate compile
