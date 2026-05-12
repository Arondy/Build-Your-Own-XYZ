package compressor

import (
	"testing"
)

func TestTreeBuild(t *testing.T) {
	freqs := []Node{{'C', 32, nil, nil}, {'D', 42, nil, nil}, {'E', 120, nil, nil}, {'K', 7, nil, nil}, {'L', 42, nil, nil}, {'M', 24, nil, nil}, {'U', 37, nil, nil}, {'Z', 2, nil, nil}}

	root := BuildHuffmanTree(freqs)
	if root.Weight != 306 {
		t.Fail()
	}
}

func TestTraverse(t *testing.T) {
	freqs := []Node{{'C', 32, nil, nil}, {'D', 42, nil, nil}, {'E', 120, nil, nil}, {'K', 7, nil, nil}, {'L', 42, nil, nil}, {'M', 24, nil, nil}, {'U', 37, nil, nil}, {'Z', 2, nil, nil}}

	root := BuildHuffmanTree(freqs)
	codes := make(map[byte]HuffmanCode, len(freqs))

	TraversePreorder(&root, codes, HuffmanCode{})
	if !(codes['E'].bytes[0] == 0 && codes['U'].bytes[0] != 4) {
		t.Fail()
	}
}
