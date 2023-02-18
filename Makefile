.DEFAULT_TARGET: all

GO_FILES := $(shell find . -type f -name '*.go')

all: completions/sundial.fish completions/sundial.zsh completions/sundial.bash completions/sundial.ps1 internal/core/cities.csv sundial

internal/core/cities.csv: scripts/makecsv.sh scripts/trim_csv.py
	scripts/makecsv.sh

sundial: $(GO_FILES) internal/core/cities.csv
	go build

completions/sundial.fish: sundial
	printf "%s" "$(./sundial completion fish)" > completions/sundial.fish

completions/sundial.zsh: sundial
	./sundial completion zsh > completions/sundial.zsh

completions/sundial.bash: sundial
	./sundial completion bash > completions/sundial.bash

completions/sundial.ps1: sundial
	./sundial completion powershell > completions/sundial.ps1

install: all
	go install .

clean:
	rm -f completions/*
	rm -f internal/core/cities.csv
	rm -f sundial

release: clean all
	go mod vendor
	go mod tidy
	goreleaser release --clean

.PHONY: all clean install release
