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

	flag.Parse()

	otherArgs := flag.Args()

	if len(otherArgs) == 0 {
		return fmt.Errorf("No file provided!")
	}

	if *countBytesFlag {
		filepath := otherArgs[0]
		bytesNumber, err := p.WC.CountBytes(filepath)

		if err != nil {
			return err
		}

		fmt.Println(bytesNumber, filepath)
	} else {
		return fmt.Errorf("Incorrect arguments provided!")
	}

	return nil
}
