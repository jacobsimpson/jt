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
    ;

selection
    : REGULAR_EXPRESSION
    | expression
    ;

expression
    : LPAREN paren=expression RPAREN
    | NOT negative=expression
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
        | DOUBLE
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

block
    : LBRACE command RBRACE
    ;

command
    : IDENTIFIER parameterList?
    ;

parameterList
    : '(' variable* ')'
    ;

variable
    : COLUMN slice?
    | IDENTIFIER slice?
    ;

slice            : '[' left=INTEGER? ':' right=INTEGER? ']' ;

IDENTIFIER       : [a-zA-Z_] [a-zA-Z_0-9]* ;
COLUMN           : '%' [0-9]+ | '%#' ;

REGULAR_EXPRESSION
    : '/' ~('/')* '/'
    | '|' ~('|')* '|'
    ;

//
// Boolean operators.
//
AND              : 'and' ;
OR               : 'or' ;
NOT              : 'not';
TRUE             : 'true' ;
FALSE            : 'false' ;
LT               : '<' ;
LE               : '<=' ;
EQ               : '==' ;
NE               : '!=' ;
GE               : '>=' ;
GT               : '>' ;

//
// Syntax tokens
//
LPAREN           : '(' ;
RPAREN           : ')' ;
LBRACE           : '{' ;
RBRACE           : '}' ;

//
// Literals
//
STRING           : '"' ~('"')* '"' ;
INTEGER          : '-'? [0-9][0-9_]* ;
HEX_INTEGER      : '-'? '0x' [0-9][0-9_]* ;
BINARY_INTEGER   : '-'? '0b' [01][01_]* ;
DOUBLE           : '-'? [0-9][0-9_]* '.' ( [0-9][0-9_]* )? ;

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

WS
    : (' ' | '\t')+ -> channel(HIDDEN)
    ;
