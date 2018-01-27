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
    | COLUMN '~' REGULAR_EXPRESSION
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
