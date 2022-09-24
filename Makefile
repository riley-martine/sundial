.DEFAULT_TARGET: all

all: static/cities.csv sundial

static/cities.csv:
	scripts/makecsv.sh

sundial: main.go
	go build

install: sundial
	rm -f ~/bin/sundial
	cp sundial ~/bin/sundial

clean:
	rm -f static/cities.csv
	rm sundial
