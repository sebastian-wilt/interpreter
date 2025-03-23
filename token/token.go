package token

type Token struct {
	Kind  TokenType
	Value string
	Pos   Position
}

type TokenType int

const (
	// Single symbols
	LEFT_PAREN    TokenType = iota // (
	RIGHT_PAREN                    // )
	LEFT_BRACE                     // {
	RIGHT_BRACE                    // }
	LEFT_BRACKET                   // [
	RIGHT_BRACKET                  // ]
	COMMA                          // ,
	SEMICOLON                      // ;
	COLON                          // :
	UNDERSCORE                     // _

	// Operators
	PLUS            // +
	PLUS_EQUAL      // +=
	MINUS           // -
	MINUS_EQUAL     // -=
	SLASH           // /
	SLASH_EQUAL     // /=
	STAR            // *
	STAR_EQUAL      // *=
	STAR_STAR       // **
	STAR_STAR_EQUAL // **=
	BANG            // !
	BANG_EQUAL      // !=
	EQUAL           // =
	EQUAL_EQUAL     // ==
	GREATER         // >
	GREATER_EQUAL   // >=
	LESS            // <
	LESS_EQUAL      // <=
	SINGLE_ARROW    // ->

	LAND        // &&
	LOR         // ||
	AND         // &
	AND_EQUAL   // &=
	OR          // |
	OR_EQUAL    // |=
	TILDE       // ~
	TILDE_EQUAL // ~=
	CARET       // ^
	CARET_EQUAL // ^=

	// Literals
	IDENT
	STRING
	CHAR
	INTEGER
	REAL

	// Keywords
	IF       // if
	ELSE     // else
	FALSE    // false
	TRUE     // true
	FOR      // for
	IN       // in
	WHILE    // while
	FUN      // fun
	RETURN   // return
	VAL      // val
	VAR      // var
	CONTINUE // continue
	FALL     // fall
	MATCH    // match

	EOF
	ILLEGAL
)
