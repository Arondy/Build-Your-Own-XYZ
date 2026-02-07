package parser_test

import (
	"json_parser/parser"
	"testing"
)

func TestStep1(t *testing.T) {
	t.Run("Right map", func(t *testing.T) {
		json := []byte("{}")
		tokens, err := parser.Lex(json)

		if err != nil {
			t.Errorf("'%s' shouldn't return lexing error: %s", json, err)
		}

		p := parser.Parser{Tokens: tokens}
		err = p.Parse()

		if err != nil {
			t.Errorf("'%s' shouldn't return parsing error: %s", json, err)
		}
	})

	t.Run("Wrong map (no closing brace)", func(t *testing.T) {
		json := []byte("{")
		tokens, err := parser.Lex(json)

		if err != nil {
			t.Errorf("'%s' shouldn't return lexing error: %s", json, err)
		}

		p := parser.Parser{Tokens: tokens}
		err = p.Parse()

		if err == nil {
			t.Errorf("'%s' should return parsing error", json)
		}
	})

	t.Run("Wrong map (no opening brace)", func(t *testing.T) {
		json := []byte("}")
		tokens, err := parser.Lex(json)

		if err != nil {
			t.Errorf("'%s' shouldn't return lexing error: %s", json, err)
		}

		p := parser.Parser{Tokens: tokens}
		err = p.Parse()

		if err == nil {
			t.Errorf("'%s' should return parsing error", json)
		}
	})

	t.Run("Wrong map (empty)", func(t *testing.T) {
		json := []byte("")
		tokens, err := parser.Lex(json)

		if err != nil {
			t.Errorf("'%s' shouldn't return lexing error: %s", json, err)
		}

		p := parser.Parser{Tokens: tokens}
		err = p.Parse()

		if err == nil {
			t.Errorf("'%s' should return parsing error", json)
		}
	})

	t.Run("Wrong map (unknown token)", func(t *testing.T) {
		json := []byte("{@")
		_, err := parser.Lex(json)

		if err == nil {
			t.Errorf("'%s' should return lexing error", json)
		}
	})
}
