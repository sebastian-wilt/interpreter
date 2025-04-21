package token

import "fmt"

type Token struct {
	Kind  TokenType
	Value string
	Pos   Position
}

// Create a new token
func NewToken(kind TokenType, value string, row int, col int) Token {
	return Token{
		Kind:  kind,
		Value: value,
		Pos: Position{
			Row:    row,
			Column: col,
		},
	}
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
	PERCENT                        // %
	UNDERSCORE                     // _

	// Operators (1-3 characters)
	PLUS            // +
	PLUS_EQUAL      // +=
	MINUS           // -
	MINUS_EQUAL     // -=
	MINUS_GREATER   // ->
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

func (t TokenType) String() string {
	switch t {
	case AND:
		return "'&'"
	case AND_EQUAL:
		return "'&='"
	case BANG:
		return "'!'"
	case BANG_EQUAL:
		return "'!="
	case CARET:
		return "'^'"
	case CARET_EQUAL:
		return "''^="
	case CHAR:
		return "char"
	case COLON:
		return "':'"
	case COMMA:
		return "','"
	case CONTINUE:
		return "'continue'"
	case ELSE:
		return "'else'"
	case EOF:
		return "''"
	case EQUAL:
		return "'='"
	case EQUAL_EQUAL:
		return "'=='"
	case FALL:
		return "'fall'"
	case FALSE:
		return "'false'"
	case FOR:
		return "'for'"
	case FUN:
		return "'fun'"
	case GREATER:
		return "'>'"
	case GREATER_EQUAL:
		return "'>='"
	case IDENT:
		return "identifier"
	case IF:
		return "'if'"
	case ILLEGAL:
		return "illegal token"
	case IN:
		return "'in'"
	case INTEGER:
		return "integer"
	case LAND:
		return "'&&'"
	case LEFT_BRACE:
		return "'{'"
	case LEFT_BRACKET:
		return "'['"
	case LEFT_PAREN:
		return "'('"
	case LESS:
		return "'<'"
	case LESS_EQUAL:
		return "'<='"
	case LOR:
		return "'||'"
	case MATCH:
		return "'match'"
	case MINUS:
		return "'-'"
	case MINUS_EQUAL:
		return "'-='"
	case MINUS_GREATER:
		return "'->'"
	case OR:
		return "'|'"
	case OR_EQUAL:
		return "'|='"
	case PLUS:
		return "'+'"
	case PLUS_EQUAL:
		return "'+='"
	case REAL:
		return "real"
	case RETURN:
		return "'return'"
	case RIGHT_BRACE:
		return "'}'"
	case RIGHT_BRACKET:
		return "']'"
	case RIGHT_PAREN:
		return "')'"
	case SEMICOLON:
		return "';'"
	case SLASH:
		return "'/'"
	case SLASH_EQUAL:
		return "'/='"
	case STAR:
		return "'*'"
	case STAR_EQUAL:
		return "'*='"
	case STAR_STAR:
		return "'**'"
	case STAR_STAR_EQUAL:
		return "'**='"
	case STRING:
		return "string"
	case TILDE:
		return "'~'"
	case TILDE_EQUAL:
		return "'~='"
	case TRUE:
		return "'true'"
	case UNDERSCORE:
		return "'_'"
	case VAL:
		return "'val'"
	case VAR:
		return "'var'"
	case WHILE:
		return "'while'"
	case PERCENT:
		return "'%"
	}

	panic(fmt.Sprintf("Unexpected token.TokenType: %#v", t))
}
