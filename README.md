# jt - Just Text

sed/awk/grep are pretty awesome tools but sometimes I wish they were slighly
different. It's odd little things, like getting substrings in awk:

```sh
substr($2, 0, length($2)-1)
```

I wish it was as easy as Python:

```sh
$2[0:-1]
```

I wish dates and times were a first class type

```sh
ps -ef | jt '%5 < 2018-01-18T06:00'
```

This is an experiment to see what that could look like. It is very nacent, and
most of the examples that follow will not actual work yet, they represent how I
would like it to work.

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
      in the jt script). If the incoming text doesn't match the type
      expectation, it isn't meeting the programmer's selection criteria.
    ```sh
    jt '%3 < "that"'
    jt '%3 < 14'
    jt '%3 < dt("12/11/2017 6:43am")'
    jt '%3 < 2017-12-11T06:43'
    jt 'dt("12/11/2017 6:00am ") < %3 < dt("12/11/2017 6:43am")'
    jt 'dt(12/11/2017 6:00am) < %3 < dt(12/11/2017 6:43am)'
    jt '%3 < 12 and %4 == "joe"'
    jt '%3 < 12 or %4 == "joe"'
    ```

- program block, colorizes the first column, prints the rest of the columns.
    ```sh
    jt '{print color(%1, blue), %[2:]}'
    ```

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

- Instead of using '$' to indicate a variable, use some other construct so that
  it's possible to cooperate easily with bash. Maybe, \1 instead of $1.
  Actually, I don't think that \1 is a good idea. \ is a bash escape character,
  so it has special significance when embedded in a string.
- Attempt to choose special characters to not conflict with bash so it's easy
  to use double quoted (") strings on the command line. That way, embedding env
  vars is easy. Hmm, maybe even make it easy to reference env vars?

- reads the text processing program from prog-file.jt
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

- simple, robust string indexing
- this will return a string that is the last 3 characters of %2, if there are 3
  characters. If there are less than 3 characters, it will return whatever it
  can, or an empty string.
    ```sh
    jt 'print %2[-3:]'
    ```

- string multiplications. "."*5 == "....."
- .lower(), .upper(), .title(), .capitalize(), .swapcase(), .reverse(),
  .join(), ltrim(), rtrim(), .trim(), .endswith(), .startswith()

- it would be nice if there was a simpler way to do this. It's such a common
  type operation for shell code to manipulate paths.
    - if we try to treat paths directly as an array, (`%PATH[0]`), which would
      be convenient, I think it will be difficult to get it right, like when to
      treat an env var as an array, and when to index a char in a string, and
      what character to use for separating the elements of the path.
    - `zsh` has a way of manipulating paths.
- and a better way of accessing environment variables would be good too. Maybe
  even just making them directly available, like `%PATH`, and `%GOPATH`.
    ```sh
    jt 'env["GOPATH"].split(":")[0]'
    ```

- remove unwanted path entries
    ```sh
    jt 'env["GOPATH"].split(":").filter("some-unwanted-path").join(":")'
    ```

- remove duplicates. Notice it won't change order.
    ```sh
    jt 'env["GOPATH"].split(":").dedup().join(":")'
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
	2012-06                  # could be an arithmetic expression
	2012-06-03               # could be an arithmetic expression, but less likely
    2012-06-03T              # I think if all the way up to the T was required
                             # for specifying a date/time literal (leaving the
                             # time parts optional), that should be sufficiently
                             # unambiguous.
	2012-06-03T23            # Unambiguously a date literal.
	d(12/11/2017 6:00am)     # Some extra work parsing this, treat d as a special
	                         # function, with special parsing rules?
	d("12/11/2017 6:00am ")  # Most regular, no special parser required, but
	                         # least usable.
	```
- add, subtract and compare dates and times. Durations should be a first class
  type too. Do best guessing to auto parse the dates and durations as part of
  type coercion. In this example, column 3 should be parsed as a number of
  different data/time formats until something works.
    ```sh
    jt '2018-12-11T06:00 < %3 < 2018-12-11T06:43'
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
    '```

- support maps, lists and sets natively.
    - make sure that the syntax for these maps onto the syntax for other,
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

I alternate between methods and functions. Not sure what to do. Methods provide
a nice opportunity for fluent programming (see the PATH split/join examples)
and functions seem a little more intuitive. If I go with methods, I don't want
null pointer crashes, so automatic elvis operator behavior.

## Development

### Preparation

- Java needs to be installed to run the antlr parser generator.
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

