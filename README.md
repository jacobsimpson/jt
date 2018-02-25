# jt - Just Text

sed/awk/grep are pretty awesome tools but sometimes I wish they were slighly
different. It's odd little things, like using substrings to remove a file
extension:

```sh
substr($2, 0, length($2)-4)
```

I wish it was as easy as Python:

```sh
$2[0:-4]
```

or maybe even easier:

```sh
%2[0:"."]
```

I wish dates and times were a first class type

```sh
ps -ef | jt '%5 < 2018-01-18T06:00'
ps -ef | jt '%5 > now - 5M'
```

Speaking of dates, haven't you always wanted to be able to use `date` not just
for producing, but also for reformatting?

```sh
ps -ef | jt 'format(%5, "%Y-%M-%DT")'
```

This is an experiment to see what that could look like. It is very nacent. See
examples of what should work in the [Cookbook](cookbook.md) doc and examples of
what I'd like to work in the [Roadmap](roadmap.md).

## Table of Contents

* [Roadmap](roadmap.md)
* [Cookbook](cookbook.md)
* [Development](#development)
    * [Preparation](#preparation)
    * [Development Cycle](#development-cycle)

## Development

### Preparation

- Java needs to be installed to run the antlr parser generator.
    - antlr parser generator jar is checked into the repo, so you don't need to
      download that.
- golang dep needs to be installed to download the necessary golang libraries.
- golang [mage](https://magefile.org/) build tool

```sh
dep ensure
mage test
```

### Development Cycle

```sh
mage
mage test
```

