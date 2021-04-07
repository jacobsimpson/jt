# JT Cookbook

Simple recipes for working with text.

## Table of Contents

- [Basic explanation](#basic-explanation)
- [Comparison operators](#comparison-operators)
- [Input column names](#input-column-names)
- [Type system](#type-system)
- [Literals](#literals)
- [Like Grep](#like-grep)
- [Matching Dates and Times](#matching-dates-and-times)

### Basic explanation

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

1. If there is nothing to match against, it is assumed that the match should be
   performed against the whole line (%0), so the example can be reduced to:

    ```sh
    jt '==/things/'
    ```

1. If there is no comparison operator, it assumed to be an equality comparison,
   so the example can be reduced to:

    ```sh
    jt '/things/'
    ```

### Comparison operators

There are the usual gang of comparison operators, `<`, `<=`, `==`, `!=`, `>=`,
and `>`.

Print all the lines where column 3 is an integer less than 3:

    jt '%3<3'

Print all the lines where the regular expression `ab[cd]` occurs somewhere in the line:

    jt '==/ab[cd]/'

#### Other valid examples

```sh
jt '<2020-01-01T'
```

- print all lines which can be coercered into a date, and are less than
  `2020-01-01T00:00:00`.

```sh
jt '>2020-02-03T'
```

- print all lines which can be coerced into a date, and are greater than
  `2020-02-03T23:59:59`.

### Input column names

This is an input stream processing language. Addressing the input is an
important part of what it does.

`jt` will produce columns for each input line.

`%[0-9]+` is used to address a particular column of input data.

`%-[0-9]+` is used to address a particular column of input data, starting from
the end of the column list and moving backwards. `%-1` is the last column of
data.

`%0` gets access to the whole line of input.

The `%` character is used to access column data to avoid conflicting with the
`$` that is widely used by Unix shells to access shell variables. That will
make commands like this easily available, which will print all rows where the
third column is numeric and greater than 3:

```sh
a=3
jt "%3 > $a"
```

### Type system

`jt` recognizes a few different types. Integer, reals, strings, dates and
regular expressions. Values that are read from data input have no type (they
are the `any` type). Input values will be coerced to match the type they are
compared again. If they can not be coerced, the match is always false.

- string
- integer
- reals
- date/time (date, timestamp, time??)
- regular expressions

There is an `any` type, which is the type of the input columns. An `any` type
means that the data hasn't yet received a type. Although this data is
represented as a string internally, during comparison operations `any` values
will behave differently from a string. Data represented as the `any` type will
be coerced to match the type it is being compared to. If it can not be coerced,
then there is no match, no matter what comparison is occurring.  `any` to `any`
comparisons are string to string comparisons.

#### Literals

There are a few types of literals supported:

-   Integers: `1`, `-10`, `0b001`, `-0xA`, `0o127`, `127_981`
-   Dates: `2012-06-01T`
-   Regular expressions: `/ab[cd]/`
-   Reals: `2.5644`
-   Strings: `"ab"`

#### Type coercion rules

`jt` will attempt to coerce to match the literal, if a literal is involved. So,
for the examples like:

```sh
jt '%3 == 2006T'
```

`%3` will be coerced to a date/time type. If `%3` can't be coerced to a
date/time, then the literal `2006T` will _not_ be coerced to a string. Instead,
the expression will evaluate to false.

More examples:

```
$ echo "1\n2\n3" | ./jt '<2006-03-01T'
$ echo "1\n2\n3" | ./jt '>=2006-03-01T'
$ echo "1\n2\n3\n2006-04-01T" | ./jt '<2006-03-01T'
$ echo "1\n2\n3\n2006-04-01T" | ./jt '>=2006-03-01T'
2006-04-01T
```

### Strings

- Substrings
    - `s = "ab.cd.txt"`
    - `s[:] == s[0:] == s`
    - `s[0:-1] == "ab.cd.tx"`

### Integers

- Integer literals can be represented in decimal, hexadecimal or binary
  formats, with `_` included for formatting purposes. Embedded `_` has no
  semantic value, it is just for visual purposes. Here are some example integer
  literals:
    - `1000 == 1_000 == 1_0_0_0`
    - `0o127`
    - `0x0001 == 0x00_01 == 0x0_0_0_0_1 == 0b1`
    - `0b1111_0001`

### Regular Expressions

- Regular expressions are a first class type with a literal representation in
  the language. The literal representation is delimited by either `/` or `|`,
  much like a string is delimited by `"` or `'` in many languages. Regular
  expressions use the extended regular expressions syntax from
  [Golang](https://golang.org/pkg/regexp/syntax/). Here are some example
  regular expression literals:
    - `/a[bc]d/`
    - `|things.|`

### Like grep

```sh
ps -ef | jt '/jt/'
```
