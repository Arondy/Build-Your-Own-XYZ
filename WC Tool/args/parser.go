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
	countWordsFlag := flag.Bool("w", true, "count words in file")

	flag.Parse()

	otherArgs := flag.Args()
	if len(otherArgs) == 0 {
		return fmt.Errorf("No file provided!")
	}

	boolFlagsNumber := 0
	flag.Visit(func(f *flag.Flag) { boolFlagsNumber += 1 })

	if boolFlagsNumber > 1 {
		return fmt.Errorf("All bool flags are mutually exclusive!")
	}

	filepath := otherArgs[0]
	file, err := p.WC.GetFileContents(filepath)
	if err != nil {
		return err
	}

	if *countBytesFlag {
		bytesNumber := p.WC.CountBytes(file)
		fmt.Println(bytesNumber, filepath)
	} else if *countLinesFlag {
		linesNumber := p.WC.CountLines(file)
		fmt.Println(linesNumber, filepath)
	} else if *countWordsFlag {
		wordsNumber := p.WC.CountWords(file)
		fmt.Println(wordsNumber, filepath)
	} else {
		return fmt.Errorf("Incorrect arguments provided!")
	}

	return nil
}
