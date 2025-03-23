package token

type Token struct {
	Kind  TokenKind
	Value string
	Row   int
	Col   int
}

type TokenKind int

const (
	// Single symbols
	LeftParen = iota
	RightParen
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket
	Comma
	Semicolon
	Colon
	Underscore

	// Operators
	Plus
	PlusEqual
	Minus
	MinusEqual
	Slash
	SlashEqual
	Star
	StarEqual
	StarStar
	StarStarEqual
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual
	SingleArrow // ->

	Land
	Lor
	And
	Or
	Not
	Caret

	// Literals
	Ident
	String
	Integer
	Real

	// Keywords
	If 
	Else 
	False 
	True 
	For 
	In 
	While 
	Fun 
	Return 
	Val 
	Var
	Continue
	Fallthrough
	Match

	EOF
	Illegal
)
