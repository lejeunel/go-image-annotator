# ====== CONFIG ======
MODULE := github.com/lejeunel/go-image-annotator

SPEC := assets/openapi.yaml
OAPI := oapi-codegen
MODELS_PKG := adapters/api/models
SERVER_PKG := adapters/api/server
MODELS_OUT := $(MODELS_PKG)/models.gen.go
SERVER_OUT := $(SERVER_PKG)/server.gen.go
VALID_AUTH_OUT := modules/auth/valid_methods.gen.go

CSS_MAIN := assets/app.css
CSS_OUT := assets/static/styles.css

# ====== TARGETS ======

.PHONY: all api-code clean

all: api-code css auth-valid-methods

auth-valid-methods: $(VALID_AUTH_OUT)
api-code: $(MODELS_OUT) $(SERVER_OUT)

css:
	tailwindcss -i $(CSS_MAIN) -o $(CSS_OUT) --minify

# --- Generate models (types only) ---
$(MODELS_OUT): $(SPEC)
	mkdir -p $(MODELS_PKG)
	$(OAPI) \
		-generate types \
		-package models \
		-o $(MODELS_OUT) \
		$(SPEC)

# --- Generate server (interfaces only, using models) ---
$(SERVER_OUT): $(SPEC) $(MODELS_OUT)
	mkdir -p $(SERVER_PKG)
	$(OAPI) \
		-generate types,std-http-server \
		-package server \
		-o $(SERVER_OUT) \
		-import-mapping $(MODULE)/$(MODELS_PKG):$(MODULE)/$(MODELS_PKG) \
		$(SPEC)


docs-dev:
	cd docs && hugo server --gc --minify --disableFastRender --logLevel debug --baseURL http://localhost:1313

$(VALID_AUTH_OUT):
	go generate ./modules/auth

# --- Cleanup generated files ---
clean:
	rm -f $(MODELS_OUT) $(SERVER_OUT) $(CSS_OUT) $(VALID_AUTH_OUT)
