package parser_test

import (
	"json_parser/parser"
	"testing"
)

type Test struct {
	name         string
	json         []byte
	wantLexErr   bool
	wantParseErr bool
}

func testCases(t *testing.T, tests map[string][]Test) {
	for groupName, group := range tests {
		t.Run(groupName, func(t *testing.T) {
			for _, tt := range group {
				t.Run(tt.name, func(t *testing.T) {
					tokens, err := parser.Lex(tt.json)

					if (err != nil) != tt.wantLexErr {
						t.Errorf("Lex() error = %v, wantLexErr %v", err, tt.wantLexErr)
					}
					if err != nil {
						return
					}

					p := parser.Parser{Tokens: tokens}
					err = p.Parse()

					if (err != nil) != tt.wantParseErr {
						t.Errorf("Parse() error = %v, wantParseErr %v", err, tt.wantParseErr)
					}
				})
			}
		})
	}
}

func TestStep1(t *testing.T) {
	tests := map[string][]Test{
		"Valid cases": {
			{name: "empty object", json: []byte("{}"), wantLexErr: false, wantParseErr: false},
		},
		"Invalid cases": {
			{name: "no closing brace", json: []byte("{"), wantLexErr: false, wantParseErr: true},
			{name: "no opening brace", json: []byte("}"), wantLexErr: false, wantParseErr: true},
			{name: "empty input", json: []byte(""), wantLexErr: false, wantParseErr: true},
			{name: "unknown token", json: []byte("{@"), wantLexErr: true, wantParseErr: false},
		},
	}

	testCases(t, tests)
}

func TestStep2(t *testing.T) {
	tests := map[string][]Test{
		"Valid cases": {
			{name: "single object", json: []byte("{\"key\": \"value\"}"), wantLexErr: false, wantParseErr: false},
			{name: "several objects", json: []byte("{\"key1\": \"value\", \"key2\": \"value\", \"key3\": \"value\"}"), wantLexErr: false, wantParseErr: false},
		},
		"Invalid cases": {
			{name: "no key", json: []byte("{: \"value\"}"), wantLexErr: false, wantParseErr: true},
			{name: "no colon", json: []byte("{\"key\" \"value\"}"), wantLexErr: false, wantParseErr: true},
			{name: "no value", json: []byte("{\"key\": }"), wantLexErr: false, wantParseErr: true},
			{name: "no quotes for string key", json: []byte("{key: \"value\"}"), wantLexErr: true, wantParseErr: false},
			{name: "no quotes for string value", json: []byte("{\"key\": value}"), wantLexErr: true, wantParseErr: false},
			{name: "token after very end", json: []byte("{} {}"), wantLexErr: false, wantParseErr: true},
			{name: "comma upfront", json: []byte("{,\"key\": \"value\"}"), wantLexErr: false, wantParseErr: true},
			{name: "comma in the very end", json: []byte("{\"key\": \"value\"},"), wantLexErr: false, wantParseErr: true},
			{name: "trailing comma after last object", json: []byte("{\"key\": \"value\",}"), wantLexErr: false, wantParseErr: true},
		},
	}

	testCases(t, tests)
}

func TestStep3(t *testing.T) {
	tests := map[string][]Test{
		"Valid cases": {
			{name: "different type values", json: []byte("{\"key1\": true,\"key2\": false,\"key3\": null,\"key4\": \"value\",\"key5\": 101}"), wantLexErr: false, wantParseErr: false},
			{name: "float number", json: []byte("{\"key\": 101.12345}"), wantLexErr: false, wantParseErr: false},
			{name: "float number starts with 0", json: []byte("{\"key\": 0.12345}"), wantLexErr: false, wantParseErr: false},
			{name: "exponential with int base number", json: []byte("{\"key\": 12e6}"), wantLexErr: false, wantParseErr: false},
			{name: "exponential with float base number", json: []byte("{\"key\": 1.2e6}"), wantLexErr: false, wantParseErr: false},
			{name: "exponential plus number", json: []byte("{\"key\": 2e+6}"), wantLexErr: false, wantParseErr: false},
			{name: "exponential minus number", json: []byte("{\"key\": 2e-6}"), wantLexErr: false, wantParseErr: false},
			{name: "exponential minus number with minus base", json: []byte("{\"key\": -2e-6}"), wantLexErr: false, wantParseErr: false},
		},
		"Invalid cases: bool, null": {
			{name: "typo in true", json: []byte("{\"key1\": truE,\"key2\": false,\"key3\": null,\"key4\": \"value\",\"key5\": 101}"), wantLexErr: true, wantParseErr: false},
			{name: "typo in true 2", json: []byte("{\"key1\": trueF,\"key2\": false,\"key3\": null,\"key4\": \"value\",\"key5\": 101}"), wantLexErr: true, wantParseErr: false},
			{name: "typo in false", json: []byte("{\"key1\": true,\"key2\": False,\"key3\": null,\"key4\": \"value\",\"key5\": 101}"), wantLexErr: true, wantParseErr: false},
			{name: "typo in false 2", json: []byte("{\"key1\": true,\"key2\": fal2se,\"key3\": null,\"key4\": \"value\",\"key5\": 101}"), wantLexErr: true, wantParseErr: false},
			{name: "typo in null", json: []byte("{\"key1\": true,\"key2\": false,\"key3\": nUll,\"key4\": \"value\",\"key5\": 101}"), wantLexErr: true, wantParseErr: false},
			{name: "typo in null 2", json: []byte("{\"key1\": true,\"key2\": false,\"key3\": nullnull,\"key4\": \"value\",\"key5\": 101}"), wantLexErr: false, wantParseErr: true},
		},
		"Invalid cases: numbers": {
			{name: "0 first digit", json: []byte("{\"key1\": true,\"key2\": false,\"key3\": null,\"key4\": \"value\",\"key5\": 01.2}"), wantLexErr: true, wantParseErr: false},
			{name: "dot in the end", json: []byte("{\"key1\": true,\"key2\": false,\"key3\": null,\"key4\": \"value\",\"key5\": 1.}"), wantLexErr: true, wantParseErr: false},
			{name: "dot upfront", json: []byte("{\"key1\": true,\"key2\": false,\"key3\": null,\"key4\": \"value\",\"key5\": .1}"), wantLexErr: true, wantParseErr: false},
			{name: "several minuses", json: []byte("{\"key1\": true,\"key2\": false,\"key3\": null,\"key4\": \"value\",\"key5\": --1.1}"), wantLexErr: true, wantParseErr: false},
			{name: "minus after dot", json: []byte("{\"key1\": true,\"key2\": false,\"key3\": null,\"key4\": \"value\",\"key5\": 1.-1}"), wantLexErr: true, wantParseErr: false},
			{name: "with float exponent", json: []byte("{\"key1\": true,\"key2\": false,\"key3\": null,\"key4\": \"value\",\"key5\": 2e1.1}"), wantLexErr: true, wantParseErr: false},
			{name: "minus in the end", json: []byte("{\"key1\": true,\"key2\": false,\"key3\": null,\"key4\": \"value\",\"key5\": 2e1-}"), wantLexErr: true, wantParseErr: false},
			{name: "unfinished exponent", json: []byte("{\"key1\": true,\"key2\": false,\"key3\": null,\"key4\": \"value\",\"key5\": 2e}"), wantLexErr: true, wantParseErr: false},
		},
	}

	testCases(t, tests)
}

func TestStep4(t *testing.T) {
	tests := map[string][]Test{
		"Valid cases": {
			{name: "different type values", json: []byte("{\"key\": \"value\",\"key-n\": 101,\"key-o\": {\"inner key\": \"inner value\"},\"key-l\": [1, 100, \"2\", null, {\"s\": 1e9}]}"), wantLexErr: false, wantParseErr: false},
			{name: "array instead of map", json: []byte("[]"), wantLexErr: false, wantParseErr: false},
		},
		"Invalid cases": {
			{name: "no closing bracket", json: []byte("["), wantLexErr: false, wantParseErr: true},
			{name: "no closing bracket 2", json: []byte("{\"key\": [}"), wantLexErr: false, wantParseErr: true},
			{name: "no opening bracket", json: []byte("]"), wantLexErr: false, wantParseErr: true},
			{name: "trailing comma", json: []byte("{\"key\": [1, 2,]}"), wantLexErr: false, wantParseErr: true},
			{name: "single quotes", json: []byte("{\"key\": 'value'}"), wantLexErr: true, wantParseErr: false},
		},
	}

	testCases(t, tests)
}
