# WARNING: Do not EDIT or MERGE this file, it is generated by 'packagespec lock'.
# build.mk builds the packages defined in packages.lock, first building all necessary
# builder images.
#
# NOTE: This file should always run as though it were in the repo root, so all paths
# are relative to the repo root.

# Include config.mk relative to repo root.
include $(shell git rev-parse --show-toplevel)/packages*.lock/config.mk

ifeq ($(PACKAGE_SPEC_ID),)
$(error You must set PACKAGE_SPEC_ID; 'make build' does this for you.)
endif

ifneq ($(PRODUCT_VERSION),)
$(error You cannot set PRODUCT_VERSION for local builds, did you mean PRODUCT_REVISION?)
endif

# PACKAGES_ROOT holds the package store, as well as other package aliases.
PACKAGES_ROOT := $(CACHE_ROOT)/packages
# PACKAGE_STORE is where we store all the package files themselves
# addressed by their input hashes.
PACKAGE_STORE := $(PACKAGES_ROOT)/store
# BY_ALIAS is where we store alias symlinks to the store.
BY_ALIAS      := $(PACKAGES_ROOT)/by-alias

# Include the layers driver.
include $(LOCKDIR)/layer.mk

# GET_IMAGE_MARKER_FILE gets the name of the Docker image marker file
# for the named build layer.
GET_IMAGE_MARKER_FILE = $($(1)_IMAGE)
# GET_IMAGE_NAME gets the Docker image name of the build layer.
GET_IMAGE_NAME = $($(1)_IMAGE_NAME)

# Determine the top-level build layer.
BUILD_LAYER_NAME      := $(shell $(call QUERY_PACKAGESPEC,.meta.builtin.BUILD_LAYERS[0].name))
BUILD_LAYER_IMAGE      = $(call GET_IMAGE_MARKER_FILE,$(BUILD_LAYER_NAME))
BUILD_LAYER_IMAGE_NAME = $(call GET_IMAGE_NAME,$(BUILD_LAYER_NAME))

BUILD_COMMAND := $(shell $(call QUERY_PACKAGESPEC,.["build-command"]))
BUILD_ENV     := $(shell $(call QUERY_PACKAGESPEC,.inputs | to_entries[] | "\(.key)=\(.value)"))
ALIASES       := $(shell $(call QUERY_PACKAGESPEC,.aliases[] | "\(.type)/\(.path)"))
ALIASES       := $(addprefix $(BY_ALIAS)/,$(ALIASES))

ifeq ($(BUILD_COMMAND),)
$(error Unable to find build command for package spec ID $(PACKAGE_SPEC_ID))
endif
ifeq ($(BUILD_ENV),)
$(error Unable to find build inputs for package spec ID $(PACKAGE_SPEC_ID))
endif

# Configure paths and filenames.
OUTPUT_DIR := $(PACKAGE_STORE)
_ := $(shell mkdir -p $(OUTPUT_DIR))
# PACKAGE_NAME is the input-addressed name of the package.
PACKAGE_NAME := $(PACKAGE_SOURCE_ID)-$(PACKAGE_SPEC_ID)
PACKAGE_ZIP_NAME := $(PACKAGE_NAME).zip
PACKAGE := $(OUTPUT_DIR)/$(PACKAGE_ZIP_NAME)
META_JSON_NAME := $(PACKAGE_ZIP_NAME).meta.json
META := $(OUTPUT_DIR)/$(META_JSON_NAME)

# In the container, place the output dir at root. This makes 'docker cp' easier.
CONTAINER_OUTPUT_DIR := /$(OUTPUT_DIR)

FULL_BUILD_COMMAND := export $(BUILD_ENV) && mkdir -p $(CONTAINER_OUTPUT_DIR) && $(BUILD_COMMAND)

### Docker run command configuration.

DOCKER_SHELL := /bin/bash -euo pipefail -c

DOCKER_RUN_ENV_FLAGS := \
	-e PACKAGE_SOURCE_ID=$(PACKAGE_SOURCE_ID) \
	-e OUTPUT_DIR=$(CONTAINER_OUTPUT_DIR) \
	-e PACKAGE_ZIP_NAME=$(PACKAGE_ZIP_NAME)

BUILD_CONTAINER_NAME := build-$(PACKAGE_SPEC_ID)-$(PACKAGE_SOURCE_ID)
DOCKER_RUN_FLAGS := $(DOCKER_RUN_ENV_FLAGS) --name $(BUILD_CONTAINER_NAME)
# DOCKER_RUN_COMMAND ties everything together to build the final package as a
# single docker run invocation.
DOCKER_RUN_COMMAND = docker run $(DOCKER_RUN_FLAGS) $(BUILD_LAYER_IMAGE_NAME) $(DOCKER_SHELL) '$(FULL_BUILD_COMMAND)'
# DOCKER_CP_COMMAND copies the built artefact from the build container.
DOCKER_CP_COMMAND = docker cp $(BUILD_CONTAINER_NAME):$(CONTAINER_OUTPUT_DIR)/$(PACKAGE_ZIP_NAME) $(PACKAGE)

# package builds the package according to the set PACKAGE_SPEC_ID and PRODUCT_REVISION.
.PHONY: package
package: $(ALIASES)
	@echo $(PACKAGE)

.PHONY: package-meta
package-meta: $(META)
	@echo $(META)

$(META): $(LOCK)
	@$(call QUERY_PACKAGESPEC,.) > $@

# PACKAGE builds the package.
$(PACKAGE): $(BUILD_LAYER_IMAGE)
	@mkdir -p $$(dirname $@)
	@echo "==> Building package: $@"
	@echo "PACKAGE_SOURCE_ID: $(PACKAGE_SOURCE_ID)"
	@echo "PACKAGE_SPEC_ID:   $(PACKAGE_SPEC_ID)"
	@# Print alias info.
	@$(call QUERY_PACKAGESPEC,.aliases[] | "alias type:\(.type) path:\(.path)") | column -t
	@docker rm -f $(BUILD_CONTAINER_NAME) > /dev/null 2>&1 || true # Speculative cleanup.
	$(DOCKER_RUN_COMMAND)
	$(DOCKER_CP_COMMAND)
	@docker rm -f $(BUILD_CONTAINER_NAME)

# ALIASES writes the package alias links.
# ALIASES must be phony to ensure they are updated to point to the
# latest builds.
.PHONY: $(ALIASES)
$(ALIASES): $(PACKAGE)
	@mkdir -p $(dir $@)
	@$(LN) -rfs $(PACKAGE) $@
	@echo "==> Package alias written: $@"
