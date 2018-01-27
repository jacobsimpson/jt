// Define a grammar called 'program' that covers the scripting
// language/expressions of the jt utility.
grammar Program;

program
    : expression*
    ;

expression
    : selection block
    | selection
    | block
    ;

selection
    : REGULAR_EXPRESSION
    | value '~' value
    ;

value
    : (
          COLUMN
        | REGULAR_EXPRESSION
        | STRING
        | DATE_TIME
        | INTEGER
        | HEX_INTEGER
        | BINARY_INTEGER
    )
    ;

STRING
    : '"' ~('"')* '"'
    ;

DATE_TIME
    : [-+]? [0-9] '-' [0-9] '-' [0-9] 'T'
    | [-+]? [0-9] '-' [0-9] '-' [0-9] 'T' [0-9] ':' [0-9] ':' [0-9]
    | [-+]? [0-9] 'T'
    | [-+]? [0-9] 'T' [0-9] ':' [0-9] ':' [0-9]
    ;

INTEGER
    : [0-9][0-9_]+
    ;

HEX_INTEGER
    : '0x' [0-9][0-9_]+
    ;

BINARY_INTEGER
    : '0b' [01][01_]+
    ;

COLUMN
    : '%' [0-9]+
    ;

block
    : '{' command* '}'
    ;

command
    : FUNCTION parameter_list?
    ;

parameter_list
    : '(' ')'
    ;

REGULAR_EXPRESSION
    : '/' ~('/')* '/'
    | '|' ~('|')* '|'
    ;

FUNCTION
    : 'print'
    ;

WS
    : (' ' | '\t')+ -> channel(HIDDEN)
    ;
