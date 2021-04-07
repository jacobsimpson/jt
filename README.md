# jt - Just Text

`sed`/`awk`/`grep` are pretty awesome tools but sometimes I wish they were
slighly different. It's odd little things, like using substrings to remove a
file extension:

```sh
ls | awk '{print substr($2, 0, length($2)-4)}'
```

I wish it was as easy as Python:

```sh
ls | jt '%2[:-4]'
```

or maybe even easier:

```sh
ls | jt '%2[:-"."]'
```

I wish often wish for just a little more date power. What if dates and times
were a first class type:

```sh
ps -ef | jt '%5 < 2018-01-18T06:00'
ps -ef | jt '%5 > now - 5M'
```

Isolating a single column of data is one of the most frequent things I do with `awk`:

```sh
ps -ef | awk '{print $2}'
```

Wouldn't it be nice if it was just a little easier:

```sh
ps -ef | jt '%2'
```

Speaking of dates, haven't you always wanted to be able to use `date` not just
for producing, but also for reformatting?

```sh
ps -ef | jt 'format(%5, "%Y-%M-%DT")'
```

This is an experiment to see what that could look like. It is very nacent. See
examples of what should work at the moment in the [Cookbook](cookbook.md) doc
and examples of what I'd like to work in the [Roadmap](roadmap.md).

## Table of Contents

* [Roadmap](roadmap.md)
* [Cookbook](cookbook.md)
* [Development](#development)
    * [Preparation](#preparation)
    * [Development Cycle](#development-cycle)

## Development

### Preparation

A couple of tools are required:

- make build tool
- [Pigeon parser generator](https://github.com/mna/pigeon)

### Development Cycle

To build the executable, just:

```sh
make
```

To run the full suite of automated tests:

```sh
make tests
```
