# Variables -------------------------------------------------------------------------------------------------------------

APPNAME     = oidc-client
REGISTRY   ?= docker.securekey.com/internal/hosting
PACKAGE     = $(APPNAME)
METAPKG     = $(PACKAGE)/backend
DATE       ?= $(shell date +%FT%T%z)
VERSION     = 0.1
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
LDFLAGS     = -X $(METAPKG).Version=$(VERSION) -X $(METAPKG).BuildDate=$(DATE) -X $(METAPKG).CommitSHA=$(COMMIT_SHA)
PKGS        = $(or $(PKG),$(shell $(GO) list ./...))
TESTPKGS    = $(shell $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
BIN         = $(CURDIR)/bin
GO          = go

V ?= 0
Q  = $(if $(filter 1,$V),,@)
M  = $(shell printf "\033[34;1m▶\033[0m")

# Targets to build Go tools --------------------------------------------------------------------------------------------

$(BIN):
	@mkdir -p $@

$(BIN)/%: | $(BIN) ; $(info $(M) building $(REPOSITORY)...)
	$Q tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get -u $(REPOSITORY) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

GOIMPORTS = $(BIN)/goimports
$(BIN)/goimports: REPOSITORY=golang.org/x/tools/cmd/goimports

GOLINT = $(BIN)/golint
$(BIN)/golint: REPOSITORY=golang.org/x/lint/golint

SWAG = $(BIN)/swag
$(BIN)/swag: REPOSITORY=github.com/swaggo/swag/cmd/swag

# Targets for our app --------------------------------------------------------------------------------------------------

.PHONY: all
all: $(BIN) fmt $(APPNAME);                                   @ ## Build backend API server

.PHONY: format
format: goimports fmt;                                        @ ## Formats code with fmt and goimports

.PHONY: $(APPNAME)
$(APPNAME): ; $(info $(M) building backend executable...)     @ ## Build backend API server
	$Q $(GO) build -ldflags '$(LDFLAGS)' -o $(BIN)/$(APPNAME) $(PACKAGE)/cmd/$(APPNAME)

.PHONY: fmt
fmt: ; $(info $(M) running gofmt...)                          @ ## Run gofmt on all source files
	$Q $(GO) fmt ./...

.PHONY: goimports
goimports: | $(GOIMPORTS) ; $(info $(M) running goimports...) @ ## Run goimports on all source files
	$Q $(GOIMPORTS) -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) running golint...)            @ ## Run golint
	$Q $(GOLINT) -set_exit_status backend/

.PHONY: swagger
swagger: | $(SWAG) ; $(info $(M) running swagger...)          @ ## Run swag to generate swagger docs
	$Q $(SWAG) init --parseDependency --parseInternal --parseDepth 2 -generalInfo backend/router.go

.PHONY: test
test: ; $(info $(M) running go test...)                       @ ## Run go unit tests
	$Q $(GO) test -v -cover $(TESTPKGS)

.PHONY: docker
docker: ; $(info $(M) building docker image...)	              @ ## Build docker image
	$Q docker build -t $(REGISTRY)/$(APPNAME):$(VERSION) .

.PHONY: clean
clean: ; $(info $(M) cleaning...)                             @ ## Cleanup everything
	@rm -rf $(BIN) $(CURDIR)/vendor $(FRONTEND)/dist $(FRONTEND)/node_modules

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
