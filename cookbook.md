# JT Cookbook

Simple recipes for working with text.

## Table of Contents

- [Basic Explanation](#basic-explanation)
- [Literals](#literals)
- [Like Grep](#like-grep)
- [Matching Dates and Times](#matching-dates-and-times)

### Basic Explanation

`jt` is a line by line stream processor, much like `awk`. It takes a boolean
statement, applies the statement to each line of incoming text and executes the
associated program block if the boolean evaluates to true.

    jt '%0 == /things/ { print(%0) }'

- `%0` represents the whole line
- strings delimited by `/` or `|` represent regular expressions
    - in the example `/things/` is the regular expression
- `{`, `}` contain the sequence of expressions to be evaluated.
    - In the example, there is only 1 expression, `print(%0)`.

This example will check if the line matches the regular expression `/things/`,
and if it does, it will print the line.

This behavior is common enough that there are a few useful defaults that apply.

1. If there is no action block, it is assumed that `{ print(%0) }` is the
   action block, so the example can be reduced to:
    ```sh
    jt '%0 == /things/'
    ```
1. If there is nothing to match against, a regular expression on it's own is
   assumed to be a comparison against the whole line, so the example can be
   reduced to:
    ```sh
    jt '/things/'
    ```

### Literals

- Dates are a first class type with a literal represenation in the language,
  just like an integer, a boolean or a string. The literal representation is
  basically ISO-8601, except that enough of the date must be included to have
  the `T` present. Here are some example date literals:
    - 2007-01-11T
    - 20070111T
    - 2018-02-12T14
    - 2018-02-12T14:02:01
- Integer literals can be represented in decimal, hexadecimal or binary
  formats, with `_` included for formatting purposes. Embedded `_` has no
  semantic value, it is just for visual purposes. Here are some example integer
  literals:
    - 1000 == 1_000 == 1_0_0_0
    - 0x0001 == 0x00_01 == 0x0_0_0_0_1 == 0b1
- Regular expressions are a first class type with a literal representation in
  the language. The literal representation is delimited by either `/` or `|`,
  much like a string is delimited by `"` or `'` in many languages. Regular
  expressions use the extended regular expressions syntax from
  [Golang](https://golang.org/pkg/regexp/syntax/). Here are some example
  regular expression literals:
    - /a[bc]d/
    - |things*|

### Like Grep

```sh
ps -ef | jt '/jt/'
```

### Matching Dates and Times

`jt` supports date/time as a native type, and has syntax support for date/time
literals

- Search for all rows where the 3rd column is a date before Dec 11, 2017 at
  6:43am local time. Since the 3rd column is being compared to a date, jt will
  attempt to parse it as a date/time, trying out various formats to see if one
  is successful. If none are successful, column 3 can not be coerced to a
  date/time, the comparison will be false.
    ```sh
    jt '%3 < 2017-12-11T06:43'
    ```

