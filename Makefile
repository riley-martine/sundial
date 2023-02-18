.DEFAULT_TARGET: all

GO_FILES := $(wildcard *.go)

all: static/cities.csv sundial

static/cities.csv: scripts/makecsv.sh scripts/trim_csv.py
	scripts/makecsv.sh

sundial: $(GO_FILES) static/cities.csv
	go build

install: all
	go install .

clean:
	rm -f static/cities.csv
	rm -f sundial

release: clean static/cities.csv
	go mod vendor
	go mod tidy
	goreleaser release

.PHONY: all clean install release
