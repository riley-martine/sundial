.DEFAULT_TARGET: all

GO_FILES := $(shell find . -type f -name '*.go')

all: internal/core/cities.csv sundial

internal/core/cities.csv: scripts/makecsv.sh scripts/trim_csv.py
	scripts/makecsv.sh

sundial: $(GO_FILES) internal/core/cities.csv
	go build

install: all
	go install .

clean:
	rm -f static/cities.csv
	rm -f sundial

release: clean internal/core/cities.csv
	go mod vendor
	go mod tidy
	goreleaser release --clean

.PHONY: all clean install release
