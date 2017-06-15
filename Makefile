SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go'| \
	grep -v vendor| \
	grep -v test| \
	grep -v mock)

all:	clean devtools build test

.PHONY: deps
deps:
	find . -name glide.yaml | while read gf;    \
	do                                          \
		pushd $$(dirname $$gf) > /dev/null; \
		glide install;                      \
		popd > /dev/null;                   \
	done

.PHONY: build
build:	deps $(SOURCES)
	go build go-inject/inject

.PHONY: mocks
mocks:	$(SOURCES)
	mockgen -destination=inject/mock/MockGraph.go go-inject/inject Graph
	mockgen -destination=inject/mock/MockProvider.go go-inject/inject Provider

.PHONY: test
test:	deps mocks
	ginkgo -r $(SOURCE_DIR)

.PHONY: clean
clean:
	go clean
	find . -name vendor | while read vendor; \
	do                                       \
		rm -rf $$vendor;                 \
	done

.PHONY: devtools
devtools:
	go get -u github.com/Masterminds/glide
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/golang/mock/mockgen
