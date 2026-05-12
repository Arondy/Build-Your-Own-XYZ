package args

import (
	"compression_tool/compressor"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const encodeExtension = ".enc"
const decodeExtension = ".txt"

func askForConfirmation(prompt string, args ...any) bool {
	fmt.Printf(prompt, args...)
	var response string
	fmt.Scanln(&response)

	return response == "y" || response == "Y" || response == ""
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Parse() error {
	inputFileFlag := flag.String("i", "", "input file")
	outputFileFlag := flag.String("o", "", "output file")
	forceOverwriteFlag := flag.Bool("f", false, "overwrite output file without confirmation")
	encodeFlag := flag.Bool("e", false, "use encode mode")
	decodeFlag := flag.Bool("d", false, "use decode mode")

	flag.Parse()

	if *encodeFlag && *decodeFlag {
		return fmt.Errorf("Only one mode can be used at a time: 'encode' or 'decode'")
	} else if !(*encodeFlag || *decodeFlag) {
		return fmt.Errorf("No mode has been provided: choose from 'encode' or 'decode' modes")
	}

	if *inputFileFlag == "" {
		return fmt.Errorf("You haven't provided input file")
	}

	if !fileExists(*inputFileFlag) {
		return fmt.Errorf("Input file doesn't exist: '%s'", *inputFileFlag)
	}

	if !*forceOverwriteFlag && fileExists(*outputFileFlag) {
		overwriteConfirmed := askForConfirmation("Output file '%s' already exists. Are you sure you want to overwrite it?\n[Y/n] ", *outputFileFlag)
		if !overwriteConfirmed {
			return nil
		}
	}

	if *outputFileFlag == "" {
		extension := filepath.Ext(*inputFileFlag)

		if *encodeFlag {
			*outputFileFlag = strings.TrimSuffix(*inputFileFlag, extension) + encodeExtension
		} else {
			*outputFileFlag = strings.TrimSuffix(*inputFileFlag, extension) + decodeExtension
		}
	}

	var err error
	if *encodeFlag {
		err = compressor.Encode(*inputFileFlag, *outputFileFlag)
	} else {
		err = compressor.Decode(*inputFileFlag, *outputFileFlag)
	}

	return err
}
