#!/usr/bin/env bash

# Make clean
rm -f cities15000.csv
rm -f cities15000.txt
rm -f cities15000.zip
rm -f cities.csv

wget http://download.geonames.org/export/dump/cities15000.zip

unzip cities15000.zip
rm cities15000.zip
mv cities15000.txt cities15000.csv
python3 trim_csv.py > static/cities.csv
rm cities15000.csv
