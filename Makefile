# ------------------------------------------------------------------------------
#  build

.PHONY: build
build:
	goreleaser build --rm-dist --single-target --snapshot

generate_k8s_deprecations:
	$(MAKE) -C tools/api-lifecycle-gen generate
