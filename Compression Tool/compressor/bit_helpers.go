package compressor

import (
	"bufio"
	"io"
)

type BitWriter struct {
	_byte byte
	count int
}

func (bw *BitWriter) writeToByte(bitstr string, writer *bufio.Writer) error {
	for _, bit := range bitstr {
		bw._byte = bw._byte << 1
		if bit == '1' {
			bw._byte = bw._byte | 1
		}
		bw.count++

		if bw.count == 8 {
			err := writer.WriteByte(bw._byte)
			if err != nil {
				return err
			}

			bw.count = 0
			bw._byte = 0
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
