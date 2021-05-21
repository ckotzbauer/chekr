
.PHONY: build
build: clear
	go build -o bin/api-lifecycle-gen

clear:
	rm -rf bin

generate: build
	./bin/api-lifecycle-gen api-repo ../../data/k8s_deprecations_generated.json
