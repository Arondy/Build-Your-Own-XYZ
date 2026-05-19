package args

import (
	"flag"
	"fmt"
	"os"
	"zip_file_cracker/cracker"
)

const encodeExtension = ".enc"
const decodeExtension = ".txt"

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Parse() error {
	inputZipFlag := flag.String("i", "", "input ZIP file")
	wordlistFileFlag := flag.String("w", "", "use wordlist")
	bruteforceFlag := flag.Bool("b", false, "use bruteforce")
	bruteforceMinFlag := flag.Int("min", 0, "bruteforce min characters")
	bruteforceMaxFlag := flag.Int("max", 0, "bruteforce max characters")

	flag.Parse()

	if *wordlistFileFlag != "" && *bruteforceFlag {
		return fmt.Errorf("Only one mode can be used at a time: 'wordlist' or 'bruteforce'")
	} else if !(*wordlistFileFlag != "" || *bruteforceFlag) {
		return fmt.Errorf("No mode has been provided: choose from 'wordlist' or 'bruteforce' modes")
	}

	if *bruteforceFlag {
		if *bruteforceMinFlag <= 0 || *bruteforceMaxFlag <= 0 {
			return fmt.Errorf("You have to provide positive 'min' and 'max' flags values")
		}
	}

	if *inputZipFlag == "" {
		return fmt.Errorf("You haven't provided input ZIP file")
	}

	if !fileExists(*inputZipFlag) {
		return fmt.Errorf("Input ZIP file doesn't exist: '%s'", *inputZipFlag)
	}

	c, err := cracker.NewCracker(*inputZipFlag)
	if err != nil {
		return err
	}

	if *wordlistFileFlag != "" {
		_, err = c.WordlistAttack(*wordlistFileFlag)
	} else {
		_, err = c.Bruteforce(*bruteforceMinFlag, *bruteforceMaxFlag)
	}

	return err
}
