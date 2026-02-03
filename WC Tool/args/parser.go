package args

import (
	"flag"
	"fmt"
	"wc_tool/wc"
)

type Parser struct {
	WC wc.WordCount
}

func (p Parser) Parse() error {
	countBytesFlag := flag.Bool("c", false, "count bytes in file")
	countLinesFlag := flag.Bool("l", false, "count lines in file")
	countWordsFlag := flag.Bool("w", false, "count words in file")
	countCharactersFlag := flag.Bool("m", false, "count characters in file")

	flag.Parse()

	otherArgs := flag.Args()
	if len(otherArgs) == 0 {
		return fmt.Errorf("No file provided!")
	}

	boolFlagsNumber := 0
	flag.Visit(func(f *flag.Flag) { boolFlagsNumber += 1 })

	if boolFlagsNumber == 0 {
		*countBytesFlag = true
		*countLinesFlag = true
		*countWordsFlag = true
	}

	filepath := otherArgs[0]
	file, err := p.WC.GetFileContents(filepath)
	if err != nil {
		return err
	}

	if *countLinesFlag {
		linesNumber := p.WC.CountLines(file)
		fmt.Print(linesNumber, " ")
	}
	if *countWordsFlag {
		wordsNumber := p.WC.CountWords(file)
		fmt.Print(wordsNumber, " ")
	}
	if *countBytesFlag {
		bytesNumber := p.WC.CountBytes(file)
		fmt.Print(bytesNumber, " ")
	}
	if *countCharactersFlag {
		charactersNumber := p.WC.CountCharacters(file)
		fmt.Print(charactersNumber, " ")
	}

	fmt.Print(filepath)

	return nil
}
