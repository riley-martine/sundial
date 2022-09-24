.DEFAULT_TARGET: all

GO_FILES := $(wildcard *.go)

all: static/cities.csv sundial

static/cities.csv: scripts/makecsv.sh scripts/trim_csv.py
	scripts/makecsv.sh

sundial: $(GO_FILES)
	go build

install: all
	go install .

clean:
	rm -f static/cities.csv
	rm sundial

.PHONY: install
