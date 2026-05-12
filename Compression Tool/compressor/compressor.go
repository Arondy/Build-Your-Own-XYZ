package compressor

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// both lower-case and upper-case
const asciiSymbolsNumber = 128
const PaddingOffset = 1

func getCharactersFrequency(reader *bufio.Reader) (freqs map[byte]int, err error) {
	freqs = make(map[byte]int, asciiSymbolsNumber)
	buffer := make([]byte, 4096)

	for {
		n, err := reader.Read(buffer)
		for i := range n {
			freqs[buffer[i]]++
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	if len(freqs) == 0 {
		return freqs, fmt.Errorf("Empty file provided")
	}

	return freqs, nil
}

func intToBytes(n int) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(n))
	return bytes
}

func writeFreqs(freqs map[byte]int, writer *bufio.Writer) error {
	err := writer.WriteByte(byte(len(freqs)))
	if err != nil {
		return err
	}
	err = writer.WriteByte(0)
	if err != nil {
		return err
	}

	for letter, freq := range freqs {
		err = writer.WriteByte(letter)
		if err != nil {
			return err
		}

		_, err = writer.Write(intToBytes(freq))
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func loadFreqs(file *bufio.Reader) (freqs map[byte]int, padding byte, err error) {
	freqs = make(map[byte]int, 0)

	lettersNumber, err := file.ReadByte()
	if err != nil {
		return nil, 0, err
	}
	padding, err = file.ReadByte()
	if err != nil {
		return nil, 0, err
	}

	for range lettersNumber {
		letter, err := file.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		freqBytes := make([]byte, 8)
		_, err = io.ReadFull(file, freqBytes)
		if err != nil {
			return nil, 0, err
		}

		freq := binary.LittleEndian.Uint64(freqBytes)
		freqs[letter] = int(freq)
	}

	return freqs, padding, nil
}

func writeEncodedContent(codes map[byte]HuffmanCode, reader *bufio.Reader, writer *bufio.Writer) (padding byte, err error) {
	bw := BitWriter{}
	readBuffer := make([]byte, 4096)

	for {
		n, err := reader.Read(readBuffer)
		for i := range n {
			code := codes[readBuffer[i]]
			err := bw.writeToByte(code, writer)
			if err != nil {
				return 0, err
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}
	}

	padding, err = bw.flush(writer)
	if err != nil {
		return 0, err
	}

	err = writer.Flush()
	return padding, err
}

func writePadding(padding byte, writer *bufio.Writer) error {
	if padding == 0 {
		return nil
	}
	err := writer.WriteByte(padding)
	if err != nil {
		return err
	}

	err = writer.Flush()
	return err
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
	)

	switch {
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func Encode(inputFilename, outputFilename string) error {
	file, err := os.Open(inputFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	freqs, err := getCharactersFrequency(reader)
	if err != nil {
		return err
	}

	nodes := freqsToNodes(freqs)
	root := BuildHuffmanTree(nodes)
	codes := make(map[byte]HuffmanCode, len(nodes))

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	reader.Reset(file)

	TraversePreorder(&root, codes, HuffmanCode{})

	outputFile, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	err = writeFreqs(freqs, writer)
	if err != nil {
		return err
	}

	padding, err := writeEncodedContent(codes, reader, writer)
	if err != nil {
		return err
	}

	_, err = outputFile.Seek(PaddingOffset, io.SeekStart)
	if err != nil {
		return err
	}

	writer.Reset(outputFile)
	err = writePadding(padding, writer)
	if err != nil {
		return err
	}

	srcStats, err := os.Stat(inputFilename)
	if err != nil {
		return err
	}
	encodedStats, err := os.Stat(outputFilename)
	if err != nil {
		return err
	}
	srcSizeFmt := formatSize(srcStats.Size())
	encSizeFmt := formatSize(encodedStats.Size())
	compression := 1 - float32(encodedStats.Size())/float32(srcStats.Size())

	fmt.Printf("File %s (%s) compressed in %s (%s), compression is %1.f%%\n", inputFilename, srcSizeFmt, outputFilename, encSizeFmt, compression*100)
	return nil
}

func decodeLetters(node *Node, root *Node, br *BitReader, writer *bufio.Writer) (*Node, error) {
	bit, err := br.Read()
	if err != nil {
		return nil, err
	}

	if !node.isLeaf() {
		if bit == 0 {
			node = node.Left
		} else {
			node = node.Right
		}
	}

	if node.isLeaf() {
		err := writer.WriteByte(node.Letter)
		if err != nil {
			return node, err
		}
		node = root
	}

	return node, err
}

func Decode(encodedFilename, outputFilename string) error {
	file, err := os.Open(encodedFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	freqs, padding, err := loadFreqs(reader)
	if err != nil {
		return err
	}

	nodes := freqsToNodes(freqs)
	root := BuildHuffmanTree(nodes)
	br, err := NewBitReader(reader, padding)
	node := &root

	outputFile, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)

	for {
		node, err = decodeLetters(node, &root, &br, writer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	fmt.Printf("File %s decoded to %s\n", encodedFilename, outputFilename)
	return nil
}
