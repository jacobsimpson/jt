# Roadmap

- implement duration literals
    - using strict ISO-8601 duration literals (e.g. P1Y3M ...) will mean that
      during parsing, duration literals could potentially be valid identifiers
      too. Requiring the P at the end instead of the beginning would give
      behavior more like the date/time literal. It would also make parsing a
      little more consistent.
- implement a `now()` function.
- implement time literals
    - right now I think time literals only work if there is a preceeding date.
- convert everything to be an expression
    - `if` statements (when they get implemented)
    - blocks - the value of the last expression in a block is the value of the
      block expression
    - the value of a block evaluation can be output to the screen, instead of
      needing an explicit `print` statement.
- implement the string functions
- enable chaining of string functions
- `|>` operator - take the output of the preceeding function and make it the
  first parameter of the following function.
- enable mathematical operations
- come up with type promotion rules
    - not sure what language I saw it in (maybe Scala) but one of them
      explicitly recognized how types could be promoted for comparisons. Ints
      promote to float/doubles for comparison, for example.
    - I think I'd like to say cross type comparison will always return false,
      unless there is a valid promotion (int -> double). The exception being
      strings from the input stream, in which case automatic parsing should be
      applied in an attempt to coerce the input stream type into the comparison
      type (when comparing to a date, attempt to parse as a date, when
      comparing to an int, attempt to parse as an int.)
- refactor to make the input stream coercion functions also available as
  explicit user parse functions.
- more comprehensive testing
    - I've been testing the individual features as I develop them, but when I
      try to use it on it's own, it's not very usable yet.
    - use instead of `grep` and `ag`. Look through my shell history for examples.
    - use instead of `awk`. Look through my shell history for examples.
    - use instead of `sed`.
    - Look on StackOverflow for other examples.
- implement variable assignment
- implement falsey/truthy behavior for other types
- implement native list/arrays
- implement native maps
- implement native sets
- implement user defined functions
- implement multiple processing rules in 1 script.
- `jt 'x/kadk/ {print(%0)}' <input>`
    - appears to succeed, gives no error message. I don't think it should do
      that.
- implement negative column addressing. %-1 will address the last column, %-2
  will address the second to last column.
- TypeScript has a cute little thing which returns an alternate value if the
  initial variable is falsy.
- elvis operator for safe chaining.

## Aspirational Examples

- these statements are all equivalent, and analogous to `awk '/this/ {print
  $0;}'`
    ```sh
    jt '/this/'
    jt '%0 == /this/'
    jt '%0 == /this/ %0'
    jt '%0 == /this/ print %0'  # Hopefully I can parse a single statement after
                                # the selection without the braces to indicate a
                                # block.
    jt '%0 == /this/ {print %0}'
    ```
- extended regular expression matcher, will automatically echo matching lines
  if there is no program block.
    ```sh
    jt '/this/'
    jt '%0 == /this/'
    ```

- extended regular expression matcher, will automatically echo lines where the
  third column matches the RE. Auto echo matching lines when there is no
  program block.
    ```sh
    jt '%3 == /that/'             # Regular expression matching on the 3rd column.
    jt '%3 == "that"'             # Exact string matching on the 3rd column.
    jt "%3 == 'that' {print %0}"  # Exact string matching on the 3rd column.
    ```

- The third column will be coerced into the matching type for each comparison
  expression and the appropriate comparison run. If the column value can not be
  coerced, it fails the comparison, no matter what the comparison is.
    - the reasoning is, the programmer supplied type information (the literal
      in the `jt` script). If the incoming text doesn't match the type
      expectation, it isn't meeting the programmer's selection criteria.
    ```sh
    jt '%3 < "that"'
    jt '%3 < 14'
    jt '%3 < 2017-12-11T06:43'
    jt '2017-12-11T06:00 < %3 < 2017-12-11T06:43'
    jt '%3 < 12 and %4 == "joe"'
    jt '%3 < 12 or %4 == "joe"'
    ```

- program block, colorizes the first column, prints the rest of the columns.
  Not sure if this is the right syntax. Ideally there would be a way of
  representing split the row, colorize 1 column and reconstitute the row.
    ```sh
    jt '{print color(%1, blue), %[2:]}'
    ```
    - Note, I had some thought that the colors would only have to be keywords
      for the color function, but what about passing the color values as
      parameters to user defined functions, which pass variables to the color
      function. No, I think they will have to be language wide keywords. Maybe
      just some basic colors as language wide keywords, some other way to
      specify the longer list of possible colors.

- as an optimization, if no one actually uses the columns, there is no need to
  split the columns.

- Ideally, the default output where no print block is specified would be to
  colorize the part of the string that matches the regexp. Even better, that
  would be representable in the syntax of the language, so that it is obvious
  that the default behavior for no print block is something from the language,
  rather than something unique.

- detect when output is redirected to a file and turn off colorization. Provide
  a command line parameter for this so that colorization can be on/off/auto, to
  allow people to specify what they want, for example when piping to less,
  which can handle the colorized output.

- echo the entire line.
    ```sh
    jt 'print %0'
    ```

- trim and print the second column.
    ```sh
    jt 'print %2.trim()'
    ```

- do an in place change to the existing file, filename.txt, changing all
  instances of 'this' to 'that'. Defaults to global (all matches on the line)
  instead of first match.
    ```sh
    jt -i 's|this|that|' filename.txt
    ```

- Avoid using '$' to indicate a variable as '$' is used by most shells. When
  the shell interacts the program passed to `jt`, things get complicated. If '$1'
  was valid in `jt` as a reference to the first column (as it is in `awk`),
  then certain programs are forced to escape shell expansion, somehow. For
  example, this `awk` program uses single quotes to address a column, which
  makes it difficult if you actually want to get the first argument to the
  shell running `awk`:

     ```
     awk '{print $1;}'
     ```

    - `\1` - possible, it's how `sed` works for back references, but it's also
      the `bash` escape character, so when used in a double quoted string, bash
      would interact with it.
    - `%1` - should be fine.
    - `_1` - would also work.

    ```
    jt "%1 == '501'"
    jt "_1 == '501'"
    ```

    And, checking to see if a column matches a shell argument would be:

    ```
    jt "%1 == $1"
    ```

    This would be interesting in the simple cases, but in the more complicated cases it wouldn't work very well:

    ```
    a="1 2 3"
    jt "%1 == $1"
    ```

- reads the text processing program from prog-file.
    ```sh
    jt -f prog-file.jt
    ```

- definitely allow matching and replacement over newlines.

- remove the last three characters of the first column and print that.
    ```sh
    jt 'print %1[:-3]'
    ```

- formatting, something analogous to ruby, where the name of the variable can
  be embedded in the formatting string, staying away from the standard C type
  mechanism found in C/Go that uses the % indicators.
    ```sh
    jt 'print "%1 --> %2 %3[:-2]"'
    ```

- easy access to environment variables.
    ```sh
    jt 'print env["PATH"].split(":")'
    ```

    -   What if `$abc` always meant an environment variable? Then, if the
        string is double quoted, the shell itself would expand the variable. If
        the string is single quoted, or otherwise escaped, the `jt` interpreter
        would map the `$abc` to an environment variable.

- and a better way of accessing environment variables would be good too. Maybe
  even just making them directly available, like `%PATH`, and `%GOPATH`.
    ```sh
    jt 'env["GOPATH"].split(":")[0]'
    ```
- simple, robust string indexing
- this will return a string that is the last 3 characters of %2, if there are 3
  characters. If there are less than 3 characters, it will return whatever it
  can, or an empty string.
    ```sh
    jt 'print %2[-3:]'
    ```

- string multiplications. `"."*5 == "....."`
- `.lower()`, `.upper()`, `.title()`, `.capitalize()`, `.swapcase()`, `.reverse()`,
  `.join()`, `ltrim()`, `rtrim()`, `.trim()`, `.endswith()`, `.startswith()`

- it would be nice if there was a simpler way to do this. It's such a common
  type operation for shell code to manipulate paths.
    - if we try to treat paths directly as an array, (`%PATH[0]`), which would
      be convenient, I think it will be difficult to get it right, like when to
      treat an env var as an array, and when to index a char in a string, and
      what character to use for separating the elements of the path.
    - `zsh` has a way of manipulating paths.

- remove unwanted path entries
    ```sh
    jt 'env["GOPATH"].split(":").filter("some-unwanted-path").join(":")'
    ```

- remove duplicates. Notice it doesn't require sorting in order to work, so it
  won't change order.
    ```sh
    jt 'env["GOPATH"].split(":").dedup().join(":")'
    jt 'env["GOPATH"].split(":").unique().join(":")'
    ```

- newline replacements.
    ```sh
    jt 's| |\n|'
    ```

- match newline in a regex.
    ```sh
    jt 's|abc\ndef|abc def|'
    ```

- merge columns 2, 3 and 4 into a single column (probably best to do that by
  preserving existing whitespace rather than re-introducing some kind of
  whitespace), and from there on, treat 2, 3 and 4 as if they were column 2.
- good for grabbing `ls -l` output, where columns 6, 7 and 8 are the time and
  date information.
    ```sh
    jt 'BEGIN merge(2,3,4); %2 < 2012-01-03T06:00'
    ```

- really easy date/time handling. Like, as a first class primitive, not
  different from string or integer. Coming up with some method of specifying a
  date that is simple and fluent would be great. ISO8601 would probably be
  good, but some additional flexibility.
    ```sh
    2012-06                  # could be an arithmetic expression. Instead, use
                             # 2012-06T.
	2012-06-03               # could be an arithmetic expression, but less likely
    2012-06-03T              # I think if all the way up to the T was required
                             # for specifying a date/time literal (leaving the
                             # time parts optional), that should be sufficiently
                             # unambiguous.
    2012T                    # Unambiguously a date literal, just the year
                             # granularity.
    2012-06T                 # Unambiguously a date literal, just the month
                             # granularity.
	2012-06-03T23            # Unambiguously a date literal.
	```
- add, subtract and compare dates and times. Durations should be a first class
  type too. Do best guessing to auto parse the dates and durations as part of
  type coercion. In this example, column 3 should be parsed as a number of
  different data/time formats until something works.
    ```sh
    jt '2018-12-11T06:00 < %3 < 2018-12-11T06:43'
    ```
- easily format dates
    ```sh
    ps -ef | jt 'format(%5, "%Y-%M-%DT")'
    ```

- allow integer representations in different bases.
- allow integer representations with _ separators
    ```sh
    jt '%1 < 0b0000_0000_0001'
    jt '%1 < 0x001'
    ```
- full complement of bitwise operators

- make networks addresses native data types.
    - IPv4, IPv6
    ```sh
    jt '%1 in 192.168.0.0/24'
    ```

- sets
    - `{}` empty set
    ```sh
    jt 'global {
        values = {'a', 'b'} # This is a set, because it uses {} and doesn't have :
    }
    %1 in values
    '
    ```

- support maps, lists and sets natively.
    - make sure that the syntax for these is analogous to the syntax for other,
      similar built in constructs
    - `{:}` empty map
    ```sh
    jt 'global {
        values = {
    	    'a': 1,
    		'b': 2,
        }
    }
    %1 in values values[%1]
    '
    ```

- great file based checking and manipulation?
   - %1.exists().isdirectory().iswritable().isreadable().isperm(0x660)

- I alternate between methods and functions. Not sure what to do. Methods
  provide a nice opportunity for fluent programming (see the PATH split/join
  examples) and functions seem a little more intuitive.
    - if method type functions were actually defined as regular functions that
      take the object as the first parameter, then the `.` operator becomes
      syntactic sugar and the user can choose either syntax as appropriate. And
      the `|>` operator makes functions even nicer.

        ```
        jt '
            function trimPeriod(s string) {...}

            /nothing/ { %1.trim().trimPeriod() }
            /things/ { %1 |> trim() |> trimPeriod() }
        ' input.txt

        ```

- No null pointer crashes, so what about automatic elvis operator behavior. I
  think Objective C (and probably Swift) have similar behaviors. And, if a
  block returns null, it's as if the selector didn't match the line. So the
  following program would match each line with `things` in it, invoke `nullify`
  on `%1`, return `null`, attempt to invoke `trim` on `null`, which will not
  actually make the call to `trim` because the object is `null`, instead it
  will skip invoking `trim` and carry the `null` result forward, making `null`
  the final result of the block, therefore the value of the block. The value of
  the block is implicitly printed, unless it is `null`, which it is in this
  case, so nothing will be printed for any lines matching the selection.

    ```
    jt '
        function nullify(s string) { null }
        /things/ { %1.nullify().trim() }
    ' input.txt
    ```

- Should process stdin when piped to, but recursive search when no files
  specified? (like ag does), instead of just freezing like grep does when it
  doesn't have any piped input or files specified?

- It would be nice if there was start and end specifiers, so that you would
  turn on certain selections when start was matched, and turn off certain
  selections when end was matched. Say there was some report type text file,
  with some header information, and a table of dates and other columns
  following the line with 'Results' going all the way to the blank line. This
  would adjust the value of the date column one day backward.

    ```
    jt '
        /Results/ start
            {%1.asDate() - 1D}
        /^$/ end
    ' input.txt
    ```

    ```
    jt '
        /Results/ -> /^$/ {
            {%1.asDate() - 1D}
        }
    ' input.txt
    ```

- It would be nice if the super common case of printing a modified line as a
  result of a match was a little more straightforward. What if the last
  statement of a block evaluates to the value of the block, and the value of
  the block gets printed? That would make these pairs of statements equivalent
  to each other.

  ```
  jt '{print(%0)}'
  jt '{%0}'
  ```

  ```
  jt '{print(%1)}'
  jt '{%1}'
  ```

  This particular `awk` program has been quite common (in my experience):


  ```
  awk '{print $1;}'
  ```

  It would be nice if it could be:

  ```
  jt '%1'
  ```
