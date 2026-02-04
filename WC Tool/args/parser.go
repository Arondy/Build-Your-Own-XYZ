package args

import (
	"flag"
	"fmt"
	"math"
	"wc_tool/wc"
)

type Parser struct {
	WC wc.WordCount
}

type Result struct {
	file       string
	lines      int
	words      int
	bytes      int
	characters int
}

func (r Result) Print(numberStringWidth int) {
	if r.lines != -1 {
		fmt.Printf("%*d ", numberStringWidth, r.lines)
	}
	if r.words != -1 {
		fmt.Printf("%*d ", numberStringWidth, r.words)
	}
	if r.bytes != -1 {
		fmt.Printf("%*d ", numberStringWidth, r.bytes)
	}
	if r.characters != -1 {
		fmt.Printf("%*d ", numberStringWidth, r.characters)
	}
	if r.file != "" {
		fmt.Print(r.file)
	}
}

func (p Parser) Parse() error {
	countBytesFlag := flag.Bool("c", false, "count bytes")
	countLinesFlag := flag.Bool("l", false, "count lines")
	countWordsFlag := flag.Bool("w", false, "count words")
	countCharactersFlag := flag.Bool("m", false, "count characters")

	flag.Parse()

	boolFlagsNumber := 0
	flag.Visit(func(f *flag.Flag) { boolFlagsNumber += 1 })

	if boolFlagsNumber == 0 {
		*countBytesFlag = true
		*countLinesFlag = true
		*countWordsFlag = true
	}

	otherArgs := flag.Args()

	if len(otherArgs) == 0 {
		otherArgs = append(otherArgs, "-")
	}

	results := make([]Result, 0, len(otherArgs))
	maxBytesNumber := 0

	for _, filepath := range otherArgs {
		file, err := p.WC.GetContents(filepath)
		if err != nil {
			return err
		}

		res := Result{filepath, -1, -1, -1, -1}
		bytesNumber := p.WC.CountBytes(file)

		if bytesNumber > maxBytesNumber {
			maxBytesNumber = bytesNumber
		}

		if *countLinesFlag {
			linesNumber := p.WC.CountLines(file)
			res.lines = linesNumber
		}
		if *countWordsFlag {
			wordsNumber := p.WC.CountWords(file)
			res.words = wordsNumber
		}
		if *countBytesFlag {
			res.bytes = bytesNumber
		}
		if *countCharactersFlag {
			charactersNumber := p.WC.CountCharacters(file)
			res.characters = charactersNumber
		}

		results = append(results, res)
	}

	numberStringWidth := int(math.Log10(float64(maxBytesNumber))) + 1

	for i, res := range results {
		res.Print(numberStringWidth)

		if i != len(results)-1 {
			fmt.Println()
		}
	}

	return nil
}
