# ====== CONFIG ======
VERSION := v0.1
REPO := github.com/lejeunel/go-image-annotator

SPEC := assets/openapi.yaml
OAPI := oapi-codegen
MODELS_PKG := adapters/api/models
SERVER_PKG := adapters/api/server
MODELS_OUT := $(MODELS_PKG)/models.gen.go
SERVER_OUT := $(SERVER_PKG)/server.gen.go
VALID_AUTH_OUT := modules/authorizer/valid_methods.gen.go
STATIC_DIR := assets/static

CSS_MAIN := assets/app.css
CSS_OUT := $(STATIC_DIR)/styles.css

.PHONY: all api-code clean build node-deps build-ci

all: api-code auth-valid-methods htmx alpine alpine-persist alpine-focus annotorious stoplight css build

node-deps:
	npm ci

build:
	go build \
		-ldflags "\
			-X '$(REPO)/globals.Version=$(VERSION)' \
			-X '$(REPO)/globals.Commit=$$(git rev-parse --short HEAD)' \
			-X '$(REPO)/globals.Date=$$(date -u +%Y-%m-%dT%H:%M:%SZ)'"

build-ci:
	$(MAKE) node-deps
	$(MAKE) build

css:
	tailwindcss -i $(CSS_MAIN) -o $(CSS_OUT) --minify

auth-valid-methods: $(VALID_AUTH_OUT)

api-code: $(MODELS_OUT) $(SERVER_OUT)


$(MODELS_OUT): $(SPEC)
	mkdir -p $(MODELS_PKG)
	$(OAPI) \
		-generate types \
		-package models \
		-o $(MODELS_OUT) \
		$(SPEC)

$(SERVER_OUT): $(SPEC) $(MODELS_OUT)
	mkdir -p $(SERVER_PKG)
	$(OAPI) \
		-generate types,std-http-server \
		-package server \
		-o $(SERVER_OUT) \
		-import-mapping $(REPO)/$(MODELS_PKG):$(REPO)/$(MODELS_PKG) \
		$(SPEC)


docs-dev:
	cd docs && hugo server --gc --minify --disableFastRender --logLevel debug --baseURL http://localhost:1313

$(VALID_AUTH_OUT):
	go generate ./modules/authorizer

htmx:
	wget https://cdn.jsdelivr.net/npm/htmx.org@2.0.10/dist/htmx.min.js -O $(STATIC_DIR)/htmx.js

alpine:
	wget https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js -O $(STATIC_DIR)/alpine.js

alpine-focus:
	wget https://cdn.jsdelivr.net/npm/@alpinejs/focus@3.x.x/dist/cdn.min.js -O $(STATIC_DIR)/alpine-focus.js

alpine-persist:
	wget https://cdn.jsdelivr.net/npm/@alpinejs/persist@3.x.x/dist/cdn.min.js -O $(STATIC_DIR)/alpine-persist.js

annotorious:
	wget https://cdn.jsdelivr.net/npm/@annotorious/annotorious@3.8.0/dist/annotorious.js -O $(STATIC_DIR)/annotorious.js
	wget https://cdn.jsdelivr.net/npm/@annotorious/annotorious@3.8.0/dist/annotorious.css -O $(STATIC_DIR)/annotorious.css

stoplight:
	wget https://unpkg.com/@stoplight/elements/web-components.min.js -O $(STATIC_DIR)/stoplight.js
	wget https://unpkg.com/@stoplight/elements/styles.min.css -O $(STATIC_DIR)/stoplight.css

clean:
	rm -f $(MODELS_OUT) $(SERVER_OUT) $(CSS_OUT) $(VALID_AUTH_OUT)
