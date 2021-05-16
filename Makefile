# ------------------------------------------------------------------------------
#  build

.PHONY: build
build:
	goreleaser build --rm-dist --single-target --snapshot
