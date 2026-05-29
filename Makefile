# ====== CONFIG ======
MODULE := github.com/lejeunel/go-image-annotator

SPEC := site/openapi.yaml
OAPI := oapi-codegen
MODELS_PKG := adapters/api/models
SERVER_PKG := adapters/api/server
MODELS_OUT := $(MODELS_PKG)/models.gen.go
SERVER_OUT := $(SERVER_PKG)/server.gen.go

CSS_MAIN := site/app.css
CSS_OUT := site/static/styles.css

# ====== TARGETS ======

.PHONY: all api-code clean

all: api-code css

api-code: $(MODELS_OUT) $(SERVER_OUT)
css: $(CSS_OUT)

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

$(CSS_OUT): $(CSS_MAIN)
	tailwindcss -i $(CSS_MAIN) -o $(CSS_OUT) --minify


# --- Cleanup generated files ---
clean:
	rm -f $(MODELS_OUT) $(SERVER_OUT) $(CSS_OUT)
