// Define a grammar called 'program' that covers the scripting
// language/expressions of the jt utility.
grammar Program;

program
    : processingRule*
    ;

// Would like to name this just 'rule' but that keyword conflicts with
// something in golang.
processingRule
    : selection block
    | selection
    | block
    ;

selection
    : REGULAR_EXPRESSION
    | expression
    ;

expression
    : LPAREN expression RPAREN
    | NOT expression
    | left=expression op=binary right=expression
    | comparison
    | boolean
    ;

comparison
    : left=value op=comparator right=value
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
        | DECIMAL
    )
    ;

comparator
    : LT | LE | EQ | NE | GE | GT
    ;

binary
    : AND | OR
    ;

boolean
    : TRUE | FALSE
    ;

AND        : 'and' ;
OR         : 'or' ;
NOT        : 'not';
TRUE       : 'true' ;
FALSE      : 'false' ;
GT         : '>' ;
GE         : '>=' ;
LT         : '<' ;
LE         : '<=' ;
EQ         : '==' ;
NE         : '!=' ;
LPAREN     : '(' ;
RPAREN     : ')' ;
DECIMAL    : '-'? [0-9][0-9_]* ( '.' [0-9][0-9_]* )? ;
IDENTIFIER : [a-zA-Z_] [a-zA-Z_0-9]* ;

STRING
    : '"' ~('"')* '"'
    ;

DATE_TIME
    : [-+]? [0-9][0-9][0-9][0-9] '-' [0-9][0-9] '-' [0-9][0-9] 'T'
    | [-+]? [0-9][0-9][0-9][0-9] '-' [0-9][0-9] '-' [0-9][0-9] 'T'
                [0-9][0-9]
    | [-+]? [0-9][0-9][0-9][0-9] '-' [0-9][0-9] '-' [0-9][0-9] 'T'
                [0-9][0-9] ':' [0-9][0-9]
    | [-+]? [0-9][0-9][0-9][0-9] '-' [0-9][0-9] '-' [0-9][0-9] 'T'
                [0-9][0-9] ':' [0-9][0-9] ':' [0-9][0-9]
    | [-+]? [0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9] 'T'
    | [-+]? [0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9] 'T'
                [0-9][0-9]
    | [-+]? [0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9] 'T'
                [0-9][0-9] ':' [0-9][0-9]
    | [-+]? [0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9] 'T'
                [0-9][0-9] ':' [0-9][0-9] ':' [0-9][0-9]
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
    | '%#'
    ;

block
    : '{' command* '}'
    ;

command
    : FUNCTION parameterList?
    ;

parameterList
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
