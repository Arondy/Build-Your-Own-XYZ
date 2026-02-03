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

	if *countBytesFlag {
		bytesNumber, err := p.WC.CountBytes(filepath)

		if err != nil {
			return err
		}

		fmt.Println(bytesNumber, filepath)
	} else if *countLinesFlag {
		linesNumber, err := p.WC.CountLines(filepath)

		if err != nil {
			return err
		}

		fmt.Println(linesNumber, filepath)
	} else {
		return fmt.Errorf("Incorrect arguments provided!")
	}

	return nil
}
