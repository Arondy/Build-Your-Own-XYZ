# Build Your Own XYZ

A collection of small programs built as part of coding challenges from [codingchallenges.fyi](https://codingchallenges.fyi/challenges/intro). Each project is an implementation of a classic tool or utility, written in Go.

## Projects

### JSON Parser

A JSON parser built from scratch with a lexer and recursive descent parser. It validates JSON syntax and provides detailed error messages with line and position information. The parser supports all JSON data types: objects, arrays, strings, numbers, booleans, and null.

### WC Tool

A clone of the Unix `wc` (word count) utility. It counts bytes, lines, words, and characters in files or from stdin. Supports the same flags as the original: `-c` (bytes), `-l` (lines), `-w` (words), and `-m` (characters). Can process multiple files and displays a total count. Output formatting matches the original, with right-aligned columns based on the largest number width.

### Compression Tool

A file compression utility based on Huffman coding. It encodes files using a Huffman tree built from character frequencies and decodes them back to their original form. Supports `-e` (encode) and `-d` (decode) modes, with `-i` for input file, `-o` for output file, and `-f` to force overwrite without confirmation. Reports compression ratio after encoding.

### Zip File Cracker

A ZIP archive password cracking tool supporting two attack modes: wordlist-based and brute-force. The wordlist attack (`-w`) reads passwords from a file and tests them concurrently using all available CPU cores. The brute-force mode (`-b`) iterates through all lowercase letter combinations within a given length range (`-min` / `-max`). Uses the `yeka/zip` library to detect encrypted ZIP entries and verify passwords.

## Running the Projects

Each project is a standalone Go module. To run a project:

```bash
cd "Project Name"
go run main.go [arguments]
```

Or compile and run the executable:

```bash
cd "Project Name"
go build -o app main.go
./app [arguments]
```

## More Challenges

More implementations from the coding challenges will be added over time.
