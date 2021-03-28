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

#### Other Valid Examples

```sh
jt '<2020-01-01T'
```

- print all lines which can be coercered into a date, and are less than
  2020-01-01T00:00:00.

```sh
jt '>2020-02-03T'
```

- print all lines which can be coerced into a date, and are greater than
  2020-02-03T23:59:59.

### Input Column Names

This is an input stream processing language. Addressing the input is an
important part of what it does.

`jt` will produce columns for each input line.

`%[0-9]+` is used to address a particular column of input data.

`%0` gets access to the whole line of input.

The `%` character is used to access column data to avoid conflicting with the
`$` that is widely used by Unix shells to access shell variables. That will
make commands like this easily available, which will print all rows where the
third column is numeric and greater than 3:

```sh
a=3
jt "%3 > $a"
```

### Literals

There are a few types of literals supported:

- string
- integer
- double
- date/time (date, timestamp, time??)
- duration
- regular expressions
- boolean

I think there should be an `any` type, which would be the type of the input
columns. An any type would be the type for any data that hasn't received a
type, making it distinct from a `string`. In this way, people would be able to
specify the type of data they know to be a string, and stronger type checking
could happen when the string type is specified. Otherwise, if unknown data
starts life as a string, you wouldn't be able to differentiate unknown typed
data from known typed data.

#### Type Coercion Rules

`jt` will attempt to coerce to match the literal, if a literal is involved. So,
for the examples like:

```sh
jt '%3 == 2006T'
```

`%3` will be coerced to a date/time type. If `%3` can't be coerced to a
date/time, then the literal `2006T` will _not_ be coerced to a string. Instead,
the expression will evaluate to false.

#### Explicit Coercions

`string(%2)`, `int(%2)`, `regex(%2)`, will attempt to coerce the second column
of input data to the specified data type. Failure to coerce will ...? (Not
sure. Something analogous to how the implicit coersion failures result in the
row not printing.)

### Complex Data Types

- set - {1, 2, 3}
- list - [1, 2, 3], ["January", "February"]
- table - {"one": 1, "two": 2}

### Strings

- Substrings
    - s = "ab.cd.txt"
    - s[:] == s[0:] == s
    - s[0:-1] == "ab.cd.tx"
    - s[:"."] == "ab"
    - s[:-"."] == "ab.cd"
    - s[".":"."] == "cd"
    - s[:/txt/] == "ab.cd."
- s.len()
- s.format("ab${c}d", {c: "3"})

### Integers

- Integer literals can be represented in decimal, hexadecimal or binary
  formats, with `_` included for formatting purposes. Embedded `_` has no
  semantic value, it is just for visual purposes. Here are some example integer
  literals:
    - 1000 == 1_000 == 1_0_0_0
    - 0o127
    - 0x0001 == 0x00_01 == 0x0_0_0_0_1 == 0b1
    - 0b1111_0001

```
1 - 3
1*3
3\2 = 1
5%3 = 2
2^3=8
3/2=1.5
```

### Regular Expressions

- Regular expressions are a first class type with a literal representation in
  the language. The literal representation is delimited by either `/` or `|`,
  much like a string is delimited by `"` or `'` in many languages. Regular
  expressions use the extended regular expressions syntax from
  [Golang](https://golang.org/pkg/regexp/syntax/). Here are some example
  regular expression literals:
    - /a[bc]d/
    - |things.|

### Dates

- Dates are a first class type with a literal represenation in the language,
  just like an integer, a boolean or a string. The complete literal
  representation is basically ISO-8601. However, there is extra support for
  partial date specifications.  Here are some example date literals:
    - 2013T
    - 2007-01-11T
    - 20070111T
    - 2018-02-12T14
    - 2018-02-12T14:02:01
- Comparisons to dates get interesting
    - Consider the date 2013T
        - `jt '<2013T'` prints all lines where the line can be coerced into a
          date, and the date is before 2013-01-01T00:00:00
        - `jt '==2013T'` prints all lines where the line can be coerced into a
          date, and the date is greater than or equal to 2013-01-01T00:00:00
          and less than 2014-01-01T00:00:00.
        - `jt '>2013T'` prints all lines where the line can be coerced into a
          date, and the date is greater than or equal to 2014-01-01T00:00:00.
        - `jt '!=2013T'` prints all lines where the line can be coerced into a
          date, and the date is less than 2013-01-01T00:00:00 or greater than
          or equal to 2014-01-01T00:00:00
- `jt` supports date/time as a native type, and has syntax support for
  date/time literals
- Print all lines where the 3rd column is a date before Dec 11, 2017 at
  6:43:00am local time. Since the 3rd column is being compared to a date, `jt`
  will attempt to parse it as a date/time, trying out various formats to see if
  one is successful. If none are successful, column 3 can not be coerced to a
  date/time, the comparison will be false.
    ```sh
    jt '%3 < 2017-12-11T06:43'
    ```

- `jt '{print 2013T-1M;}'`
- `jt '%3 < today()' - print all the lines where the 3rd column can be coerced
  into a date which is before 00:00am of today.
    - `yesterday()`
    - `tomorrow()`

#### Durations

- A duration literal starts with a P. P[n]Y[n]M[n]DT[n]H[n]M[n]S or P[n]W

### Like Grep

```sh
ps -ef | jt '/jt/'
```

### Blocks of Execution

Expressions can be executed for each matching line by enclosing them in braces:

    jt '/abc/{%3 = %3 + 10; print(%0)}'

- matches each line that has 'abc' somewhere in it. 

### Formatting

    jt '{%3 = %3.blue(); print(%0)}'

- matches all lines, changes the 3rd column to blue and prints the whole line

- There are a number of color convenience functions available.
    - `blue`
    - `red`
- Background color is available by prefacing the color name with `on_`
    - `on_blue`
- Formatting functions:
    - `italic`
    - `bold`
    - `clear`

### Printing

By default, the final expression evaluated in a block will be printed.

    jt '{print %0;}'

Can be represented more succinctly as:

    jt '{%0
