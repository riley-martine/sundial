# Sundial ☉/☾

A simple CLI program to print the percent through the day or night.
Goes well in a \*line (tmuxline, airline, powerline, etc).

```shell
☿ sundial -city Denver
35% ☉
```

## Installation

See [releases](https://github.com/riley-martine/sundial/releases), or

```shell
☿ go install github.com/riley-martine/sundial@latest
```

### Manual Installation

Requirements for building:

- Python 3
- bash
- wget
- unzip
- make
- go 1.19

```shell
☿ git clone https://github.com/riley-martine/sundial
☿ make
☿ make install
☿ sundial -version
```

## Usage

```shell
☿ sundial -city Denver
38% ☉

# International support:
# (supports all cities worldwide with population >= 15,000)
☿ sundial -city Jakarta
34% ☾

# Progressive narrowing:
☿ sundial -city Washington
could not narrow down between cities: [{Washington, GB, ENG, 54.90, -1.52} {Washington, US, DC, 38.90, -77.04} {Washington, US, IL, 40.70, -89.41} {Washington, US, UT, 37.13, -113.51}]
You may need to be more specific about which city you're in. Try specifying a country code (second field) and a fips code (third field).
        e.g. sundial -city Washington -country US -fipscode IL

☿ sundial -city Washington -country US
could not narrow down between cities: [{Washington, US, DC, 38.90, -77.04} {Washington, US, IL, 40.70, -89.41} {Washington, US, UT, 37.13, -113.51}]
You may need to be more specific about which city you're in. Try specifying a country code (second field) and a fips code (third field).
        e.g. sundial -city Washington -country US -fipscode IL

☿ sundial -city Washington -country US -fipscode IL
48% ☉

# Help text:
☿ sundial --help
Usage of sundial:
  -city string
        Name of city you're in. Required.
  -country string
        Two-letter country code, e.g. 'US'. Not required if only one city by name.
  -fipscode string
        Fipscode of region you're in. In the US, this is the two-letter state abbreviation. Otherwise, search http://download.geonames.org/export/dump/admin1CodesASCII.txt for '$countryCode.' and select the value after the period for the region you're in. Not required if only one city in country with name.
   -version
        Print program version
```

## Motivation

I think that we're too rigid in the ways we measure and think about time. This
is an experiment with a different, older way: by the passage of the sun across
the sky.

- "We're halfway to sunrise" means something very important to me that I'm not
  able to easily express with the standard 12-hour clock.

- "The days are shorter in the winter" _actually_ means that the unit "day" is
  shorter.

- The measure of time is aligned to the location. It is important to me to be
  linked to the cycles of the world, rather than linked to abstractions of them.
  (Yes, this is still an abstraction. No, you did not understand or make a good
  point.)

- It is not very good for precise coordination of meetings across timezones (for
  that, you would want something like "at exactly 2:15pm EST"). That is the
  point.

This is all an experiment. How does it _feel_ to keep time this way? What other
ways could we keep time, if we wanted? Do the ways we use fit us well, as
embodied beings? If not, what can we do about it? The full force of your
creativity is behind you!

To me, 5 months into trying this, it has been very worthwhile. It feels similar
to keeping time by the moon phase, in that there is a _real thing_ that time is
in reference to. Standard clock time reminds me of alarms that go off at exactly
the same time, while this method feels like waking up with the body's own
rhythm. It feels good to work until sunset, and I feel better about putting down
the computer, compared to if I'm thinking in terms of hours.

## Releasing

- Make changes as desired. Commit. Push.
- `git tag vX.Y.X && git push --tags`
- [GitHub actions](https://github.com/riley-martine/sundial/actions) handles the
  rest.

<!-- ### Manual releasing -->
<!-- - Install `goreleaser` ([install docs](https://goreleaser.com/install/)). -->
<!-- - `git tag vX.Y.X && git push --tags` -->
<!-- - Set `GITHUB_TOKEN` to a token with `write:packages` -->
<!-- - Run `make release`. -->
