package compressor

import (
	"bufio"
	"io"
)

type BitWriter struct {
	_byte byte
	count int
}

func (bw *BitWriter) writeToByte(code HuffmanCode, writer *bufio.Writer) error {
	for ind, _byte := range code.bytes {
		bitsCount := 8

		if ind == len(code.bytes)-1 {
			bitsCount = code.lastBitsCount
			if bitsCount == 0 {
				continue
			}
		}

		if bw.count+bitsCount <= 8 {
			src := bw._byte << bitsCount
			add := _byte >> (8 - bitsCount)
			bw._byte = src | add
			bw.count += bitsCount

			if bw.count == 8 {
				if err := writer.WriteByte(bw._byte); err != nil {
					return err
				}
				bw.count = 0
			}
		} else {
			canFit := 8 - bw.count
			src := bw._byte << canFit
			add := _byte >> (8 - canFit)

			out := src | add
			if err := writer.WriteByte(out); err != nil {
				return err
			}

			remBits := bitsCount - canFit
			var mask byte = (1 << remBits) - 1
			bw._byte = (_byte >> (8 - bitsCount)) & mask
			bw.count = remBits
		}
	}
	return nil
}

func (bw *BitWriter) flush(writer *bufio.Writer) (padding byte, err error) {
	if bw.count != 0 {
		bw._byte = bw._byte << (8 - bw.count)
		err = writer.WriteByte(bw._byte)
		return byte(8 - bw.count), err
	}
	return 0, nil
}

type BitReader struct {
	reader     *bufio.Reader
	padding    byte
	_byte      byte
	nextByte   byte
	count      byte
	isLastByte bool
}

func NewBitReader(reader *bufio.Reader, padding byte) (br BitReader, err error) {
	br = BitReader{reader: reader, padding: padding}
	br.nextByte, err = br.reader.ReadByte()

	return br, err
}

func (br *BitReader) Read() (byte, error) {
	var err error

	if br.count == 0 {
		br._byte = br.nextByte
		br.nextByte, err = br.reader.ReadByte()
		br.count = 8

		if err == io.EOF {
			if br.isLastByte {
				return 0, err
			}

			br._byte = br._byte >> br.padding
			br.count -= br.padding
			br.isLastByte = true
		} else if err != nil {
			return 0, err
		}
	}

	br.count--
	bit := (br._byte >> br.count) & 1

	return bit, nil
}
