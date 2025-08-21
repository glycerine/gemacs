package main

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/glycerine/tcell_old_hacked_up/termbox"
)

// Token types for syntax highlighting
type TokenType int

const (
	TokenNone TokenType = iota
	TokenKeyword
	TokenString
	TokenComment
	TokenNumber
	TokenOperator
	TokenIdentifier
	TokenType_
	TokenFunction
	TokenConstant
	TokenPreprocessor
	TokenSpecial
)

// Token represents a syntax element
type Token struct {
	Type   TokenType
	Start  int // byte offset in line
	End    int // byte offset in line
	Value  string
}

// Language represents a programming language definition
type Language struct {
	Name         string
	Extensions   []string
	Keywords     map[string]TokenType
	Patterns     []Pattern
	SingleComment string
	MultiComment  [2]string
	StringDelims  []string
}

// Pattern represents a regex pattern for tokenization
type Pattern struct {
	Regex *regexp.Regexp
	Type  TokenType
}

// SyntaxHighlighter manages syntax highlighting for different languages
type SyntaxHighlighter struct {
	languages map[string]*Language
	themes    map[string]*Theme
	enabled   bool
}

// Theme defines color mappings for token types
type Theme struct {
	Name   string
	Colors map[TokenType]termbox.Attribute
}

// NewSyntaxHighlighter creates a new syntax highlighter
func NewSyntaxHighlighter() *SyntaxHighlighter {
	sh := &SyntaxHighlighter{
		languages: make(map[string]*Language),
		themes:    make(map[string]*Theme),
		enabled:   true,
	}
	
	sh.initLanguages()
	sh.initThemes()
	
	return sh
}

// initLanguages initializes built-in language definitions
func (sh *SyntaxHighlighter) initLanguages() {
	// Go language definition
	go_keywords := map[string]TokenType{
		"package": TokenKeyword, "import": TokenKeyword, "func": TokenKeyword,
		"var": TokenKeyword, "const": TokenKeyword, "type": TokenKeyword,
		"struct": TokenKeyword, "interface": TokenKeyword, "map": TokenKeyword,
		"chan": TokenKeyword, "select": TokenKeyword, "case": TokenKeyword,
		"default": TokenKeyword, "if": TokenKeyword, "else": TokenKeyword,
		"for": TokenKeyword, "range": TokenKeyword, "switch": TokenKeyword,
		"goto": TokenKeyword, "break": TokenKeyword, "continue": TokenKeyword,
		"return": TokenKeyword, "defer": TokenKeyword, "go": TokenKeyword,
		"fallthrough": TokenKeyword,
		// Built-in types
		"bool": TokenType_, "byte": TokenType_, "complex64": TokenType_,
		"complex128": TokenType_, "error": TokenType_, "float32": TokenType_,
		"float64": TokenType_, "int": TokenType_, "int8": TokenType_,
		"int16": TokenType_, "int32": TokenType_, "int64": TokenType_,
		"rune": TokenType_, "string": TokenType_, "uint": TokenType_,
		"uint8": TokenType_, "uint16": TokenType_, "uint32": TokenType_,
		"uint64": TokenType_, "uintptr": TokenType_,
		// Built-in functions
		"append": TokenFunction, "cap": TokenFunction, "close": TokenFunction,
		"complex": TokenFunction, "copy": TokenFunction, "delete": TokenFunction,
		"imag": TokenFunction, "len": TokenFunction, "make": TokenFunction,
		"new": TokenFunction, "panic": TokenFunction, "print": TokenFunction,
		"println": TokenFunction, "real": TokenFunction, "recover": TokenFunction,
		// Constants
		"true": TokenConstant, "false": TokenConstant, "nil": TokenConstant,
		"iota": TokenConstant,
	}
	
	go_patterns := []Pattern{
		{regexp.MustCompile(`\b\d+(\.\d+)?([eE][+-]?\d+)?\b`), TokenNumber},
		{regexp.MustCompile(`[+\-*/%=!<>&|^~]+`), TokenOperator},
		{regexp.MustCompile(`[(){}[\],;:.]+`), TokenSpecial},
	}
	
	sh.languages["go"] = &Language{
		Name:          "Go",
		Extensions:    []string{".go"},
		Keywords:      go_keywords,
		Patterns:      go_patterns,
		SingleComment: "//",
		MultiComment:  [2]string{"/*", "*/"},
		StringDelims:  []string{`"`, "`"},
	}
	
	// JavaScript language definition
	js_keywords := map[string]TokenType{
		"var": TokenKeyword, "let": TokenKeyword, "const": TokenKeyword,
		"function": TokenKeyword, "return": TokenKeyword, "if": TokenKeyword,
		"else": TokenKeyword, "for": TokenKeyword, "while": TokenKeyword,
		"do": TokenKeyword, "switch": TokenKeyword, "case": TokenKeyword,
		"default": TokenKeyword, "break": TokenKeyword, "continue": TokenKeyword,
		"try": TokenKeyword, "catch": TokenKeyword, "finally": TokenKeyword,
		"throw": TokenKeyword, "new": TokenKeyword, "typeof": TokenKeyword,
		"instanceof": TokenKeyword, "in": TokenKeyword, "class": TokenKeyword,
		"extends": TokenKeyword, "import": TokenKeyword, "export": TokenKeyword,
		"async": TokenKeyword, "await": TokenKeyword,
		// Constants
		"true": TokenConstant, "false": TokenConstant, "null": TokenConstant,
		"undefined": TokenConstant,
	}
	
	js_patterns := []Pattern{
		{regexp.MustCompile(`\b\d+(\.\d+)?([eE][+-]?\d+)?\b`), TokenNumber},
		{regexp.MustCompile(`[+\-*/%=!<>&|^~?:]+`), TokenOperator},
		{regexp.MustCompile(`[(){}[\],;.]+`), TokenSpecial},
	}
	
	sh.languages["javascript"] = &Language{
		Name:          "JavaScript",
		Extensions:    []string{".js", ".jsx", ".mjs"},
		Keywords:      js_keywords,
		Patterns:      js_patterns,
		SingleComment: "//",
		MultiComment:  [2]string{"/*", "*/"},
		StringDelims:  []string{`"`, `'`, "`"},
	}
}

// initThemes initializes built-in color themes
func (sh *SyntaxHighlighter) initThemes() {
	// Default theme
	defaultTheme := &Theme{
		Name: "default",
		Colors: map[TokenType]termbox.Attribute{
			TokenKeyword:      termbox.ColorMagenta,
			TokenString:       termbox.ColorGreen,
			TokenComment:      termbox.ColorCyan,
			TokenNumber:       termbox.ColorYellow,
			TokenOperator:     termbox.ColorRed,
			TokenIdentifier:   termbox.ColorDefault,
			TokenType_:        termbox.ColorBlue,
			TokenFunction:     termbox.ColorBlue,
			TokenConstant:     termbox.ColorYellow,
			TokenPreprocessor: termbox.ColorMagenta,
			TokenSpecial:      termbox.ColorRed,
		},
	}
	sh.themes["default"] = defaultTheme
	
	// Dark theme
	darkTheme := &Theme{
		Name: "dark",
		Colors: map[TokenType]termbox.Attribute{
			TokenKeyword:      termbox.ColorMagenta | termbox.AttrBold,
			TokenString:       termbox.ColorGreen,
			TokenComment:      termbox.ColorCyan,
			TokenNumber:       termbox.ColorYellow,
			TokenOperator:     termbox.ColorWhite,
			TokenIdentifier:   termbox.ColorWhite,
			TokenType_:        termbox.ColorBlue | termbox.AttrBold,
			TokenFunction:     termbox.ColorCyan | termbox.AttrBold,
			TokenConstant:     termbox.ColorYellow | termbox.AttrBold,
			TokenPreprocessor: termbox.ColorMagenta,
			TokenSpecial:      termbox.ColorWhite,
		},
	}
	sh.themes["dark"] = darkTheme
}

// DetectLanguage detects language from file path
func (sh *SyntaxHighlighter) DetectLanguage(path string) *Language {
	if path == "" {
		return nil
	}
	
	ext := strings.ToLower(filepath.Ext(path))
	for _, lang := range sh.languages {
		for _, langExt := range lang.Extensions {
			if ext == langExt {
				return lang
			}
		}
	}
	
	return nil
}

// TokenizeLine tokenizes a single line of text
func (sh *SyntaxHighlighter) TokenizeLine(line []byte, lang *Language) []Token {
	if !sh.enabled || lang == nil {
		return nil
	}
	
	var tokens []Token
	text := string(line)
	pos := 0
	
	for pos < len(text) {
		// Skip whitespace
		if unicode.IsSpace(rune(text[pos])) {
			pos++
			continue
		}
		
		// Check for comments
		if token := sh.tryComment(text, pos, lang); token != nil {
			tokens = append(tokens, *token)
			pos = token.End
			continue
		}
		
		// Check for strings
		if token := sh.tryString(text, pos, lang); token != nil {
			tokens = append(tokens, *token)
			pos = token.End
			continue
		}
		
		// Check for numbers
		if token := sh.tryNumber(text, pos); token != nil {
			tokens = append(tokens, *token)
			pos = token.End
			continue
		}
		
		// Check for keywords and identifiers
		if token := sh.tryKeywordOrIdentifier(text, pos, lang); token != nil {
			tokens = append(tokens, *token)
			pos = token.End
			continue
		}
		
		// Check for operators and special characters
		if token := sh.tryOperatorOrSpecial(text, pos, lang); token != nil {
			tokens = append(tokens, *token)
			pos = token.End
			continue
		}
		
		// Skip unknown character
		pos++
	}
	
	return tokens
}

// tryComment attempts to parse a comment
func (sh *SyntaxHighlighter) tryComment(text string, pos int, lang *Language) *Token {
	// Single line comment
	if lang.SingleComment != "" && strings.HasPrefix(text[pos:], lang.SingleComment) {
		return &Token{
			Type:  TokenComment,
			Start: pos,
			End:   len(text),
			Value: text[pos:],
		}
	}
	
	// Multi-line comment (simplified - assumes it ends on same line)
	if lang.MultiComment[0] != "" && strings.HasPrefix(text[pos:], lang.MultiComment[0]) {
		start := pos
		pos += len(lang.MultiComment[0])
		end := strings.Index(text[pos:], lang.MultiComment[1])
		if end == -1 {
			end = len(text)
		} else {
			end = pos + end + len(lang.MultiComment[1])
		}
		
		return &Token{
			Type:  TokenComment,
			Start: start,
			End:   end,
			Value: text[start:end],
		}
	}
	
	return nil
}

// tryString attempts to parse a string
func (sh *SyntaxHighlighter) tryString(text string, pos int, lang *Language) *Token {
	for _, delim := range lang.StringDelims {
		if strings.HasPrefix(text[pos:], delim) {
			start := pos
			pos += len(delim)
			
			// Find closing delimiter
			for pos < len(text) {
				if strings.HasPrefix(text[pos:], delim) {
					pos += len(delim)
					break
				}
				if text[pos] == '\\' && pos+1 < len(text) {
					pos += 2 // Skip escaped character
				} else {
					pos++
				}
			}
			
			return &Token{
				Type:  TokenString,
				Start: start,
				End:   pos,
				Value: text[start:pos],
			}
		}
	}
	
	return nil
}

// tryNumber attempts to parse a number
func (sh *SyntaxHighlighter) tryNumber(text string, pos int) *Token {
	if !unicode.IsDigit(rune(text[pos])) {
		return nil
	}
	
	start := pos
	for pos < len(text) && (unicode.IsDigit(rune(text[pos])) || text[pos] == '.') {
		pos++
	}
	
	return &Token{
		Type:  TokenNumber,
		Start: start,
		End:   pos,
		Value: text[start:pos],
	}
}

// tryKeywordOrIdentifier attempts to parse a keyword or identifier
func (sh *SyntaxHighlighter) tryKeywordOrIdentifier(text string, pos int, lang *Language) *Token {
	if !unicode.IsLetter(rune(text[pos])) && text[pos] != '_' {
		return nil
	}
	
	start := pos
	for pos < len(text) && (unicode.IsLetter(rune(text[pos])) || unicode.IsDigit(rune(text[pos])) || text[pos] == '_') {
		pos++
	}
	
	word := text[start:pos]
	tokenType := TokenIdentifier
	
	if langTokenType, exists := lang.Keywords[word]; exists {
		tokenType = langTokenType
	}
	
	return &Token{
		Type:  tokenType,
		Start: start,
		End:   pos,
		Value: word,
	}
}

// tryOperatorOrSpecial attempts to parse operators or special characters
func (sh *SyntaxHighlighter) tryOperatorOrSpecial(text string, pos int, lang *Language) *Token {
	for _, pattern := range lang.Patterns {
		if pattern.Regex.MatchString(text[pos:]) {
			match := pattern.Regex.FindString(text[pos:])
			if len(match) > 0 {
				return &Token{
					Type:  pattern.Type,
					Start: pos,
					End:   pos + len(match),
					Value: match,
				}
			}
		}
	}
	
	return nil
}

// GetTheme returns the specified theme
func (sh *SyntaxHighlighter) GetTheme(name string) *Theme {
	if theme, exists := sh.themes[name]; exists {
		return theme
	}
	return sh.themes["default"]
}

// SetEnabled enables or disables syntax highlighting
func (sh *SyntaxHighlighter) SetEnabled(enabled bool) {
	sh.enabled = enabled
}

// IsEnabled returns whether syntax highlighting is enabled
func (sh *SyntaxHighlighter) IsEnabled() bool {
	return sh.enabled
}