#!/usr/bin/env bash

TMP=$(mktemp -d)

# Change zip download to change what cities are taken in.
# http://download.geonames.org/export/dump/
# Using 5000 pop gives about an extra meg to the binary, and 3ms to execution, over 150000.
wget http://download.geonames.org/export/dump/cities15000.zip -O "$TMP"/cities.zip

mkdir -p static
./scripts/trim_csv.py <(unzip -p "$TMP"/cities.zip) > static/cities.csv
rm -rf "$TMP"
