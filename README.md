# Build Your Own XYZ

A collection of small programs built as part of coding challenges from [codingchallenges.fyi](https://codingchallenges.fyi/challenges/intro). Each project is an implementation of a classic tool or utility, written in Go.

## Projects

### JSON Parser

A JSON parser built from scratch with a lexer and recursive descent parser. It validates JSON syntax and provides detailed error messages with line and position information. The parser supports all JSON data types: objects, arrays, strings, numbers, booleans, and null.

### WC Tool

A clone of the Unix `wc` (word count) utility. It counts bytes, lines, words, and characters in files or from stdin. Supports the same flags as the original: `-c` (bytes), `-l` (lines), `-w` (words), and `-m` (characters). Can process multiple files and displays a total count. Output formatting matches the original, with right-aligned columns based on the largest number width.

## Running the Projects

Each project is a standalone Go module. To run a project:

```bash
cd "Project Name"
go run main.go [arguments]
```

## More Challenges

More implementations from the coding challenges will be added over time.
