package main

import (
	"testing"
)

func TestGoTokenization(t *testing.T) {
	sh := NewSyntaxHighlighter()
	goLang := sh.languages["go"]
	if goLang == nil {
		t.Fatal("Go language not found")
	}
	
	// Test Go code tokenization
	testLine := []byte(`func main() { fmt.Println("Hello, world!") }`)
	tokens := sh.TokenizeLine(testLine, goLang)
	
	if len(tokens) == 0 {
		t.Fatal("No tokens generated")
	}
	
	// Check that we have expected tokens
	foundKeyword := false
	foundString := false
	foundIdentifier := false
	
	for _, token := range tokens {
		switch token.Type {
		case TokenKeyword:
			if token.Value == "func" {
				foundKeyword = true
			}
		case TokenString:
			if token.Value == `"Hello, world!"` {
				foundString = true
			}
		case TokenIdentifier:
			if token.Value == "main" || token.Value == "Println" {
				foundIdentifier = true
			}
		}
	}
	
	if !foundKeyword {
		t.Error("Expected to find 'func' keyword")
	}
	if !foundString {
		t.Error("Expected to find string literal")
	}
	if !foundIdentifier {
		t.Error("Expected to find identifiers")
	}
}

func TestJavaScriptTokenization(t *testing.T) {
	sh := NewSyntaxHighlighter()
	jsLang := sh.languages["javascript"]
	if jsLang == nil {
		t.Fatal("JavaScript language not found")
	}
	
	// Test JavaScript code tokenization
	testLine := []byte(`function greet(name) { console.log("Hello " + name); }`)
	tokens := sh.TokenizeLine(testLine, jsLang)
	
	if len(tokens) == 0 {
		t.Fatal("No tokens generated")
	}
	
	// Check for function keyword
	foundFunction := false
	for _, token := range tokens {
		if token.Type == TokenKeyword && token.Value == "function" {
			foundFunction = true
			break
		}
	}
	
	if !foundFunction {
		t.Error("Expected to find 'function' keyword")
	}
}

func TestLanguageDetection(t *testing.T) {
	sh := NewSyntaxHighlighter()
	
	// Test Go file detection
	goLang := sh.DetectLanguage("test.go")
	if goLang == nil || goLang.Name != "Go" {
		t.Error("Failed to detect Go language from .go extension")
	}
	
	// Test JavaScript file detection
	jsLang := sh.DetectLanguage("test.js")
	if jsLang == nil || jsLang.Name != "JavaScript" {
		t.Error("Failed to detect JavaScript language from .js extension")
	}
	
	// Test unknown extension
	unknownLang := sh.DetectLanguage("test.unknown")
	if unknownLang != nil {
		t.Error("Should not detect language for unknown extension")
	}
}

func TestThemes(t *testing.T) {
	sh := NewSyntaxHighlighter()
	
	// Test default theme
	defaultTheme := sh.GetTheme("default")
	if defaultTheme == nil {
		t.Fatal("Default theme not found")
	}
	if defaultTheme.Name != "default" {
		t.Error("Default theme has wrong name")
	}
	
	// Test dark theme
	darkTheme := sh.GetTheme("dark")
	if darkTheme == nil {
		t.Fatal("Dark theme not found")
	}
	if darkTheme.Name != "dark" {
		t.Error("Dark theme has wrong name")
	}
	
	// Test non-existent theme (should return default)
	unknownTheme := sh.GetTheme("nonexistent")
	if unknownTheme == nil || unknownTheme.Name != "default" {
		t.Error("Should return default theme for unknown theme name")
	}
}

func TestSyntaxHighlighterToggle(t *testing.T) {
	sh := NewSyntaxHighlighter()
	
	// Should be enabled by default
	if !sh.IsEnabled() {
		t.Error("Syntax highlighter should be enabled by default")
	}
	
	// Test disabling
	sh.SetEnabled(false)
	if sh.IsEnabled() {
		t.Error("Syntax highlighter should be disabled")
	}
	
	// Test re-enabling
	sh.SetEnabled(true)
	if !sh.IsEnabled() {
		t.Error("Syntax highlighter should be enabled again")
	}
}