package main

import (
	"fmt"
)

func debugTokenization() {
	sh := NewSyntaxHighlighter()
	goLang := sh.languages["go"]
	
	testLine := []byte(`func main() { fmt.Println("Hello, world!") }`)
	tokens := sh.TokenizeLine(testLine, goLang)
	
	fmt.Printf("Tokenizing: %s\n", string(testLine))
	fmt.Printf("Generated %d tokens:\n", len(tokens))
	
	for i, token := range tokens {
		var typeName string
		switch token.Type {
		case TokenKeyword:
			typeName = "Keyword"
		case TokenString:
			typeName = "String"
		case TokenIdentifier:
			typeName = "Identifier"
		case TokenSpecial:
			typeName = "Special"
		case TokenOperator:
			typeName = "Operator"
		default:
			typeName = "Other"
		}
		fmt.Printf("  Token %d: %s '%s' (%d-%d)\n", i, typeName, token.Value, token.Start, token.End)
	}
}