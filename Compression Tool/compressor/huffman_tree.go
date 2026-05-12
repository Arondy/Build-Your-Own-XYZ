package compressor

import (
	"slices"
)

type HuffmanCode struct {
	bytes         []byte
	lastBitsCount int
}

func (hc HuffmanCode) addBit(bit byte) HuffmanCode {
	if hc.lastBitsCount == 8 || len(hc.bytes) == 0 {
		hc.bytes = append(hc.bytes, 0)
		hc.lastBitsCount = 0
	} else {
		newBytes := make([]byte, len(hc.bytes))
		copy(newBytes, hc.bytes)
		hc.bytes = newBytes
	}

	hc.bytes[len(hc.bytes)-1] = hc.bytes[len(hc.bytes)-1] << 1

	if bit == 1 {
		hc.bytes[len(hc.bytes)-1] = hc.bytes[len(hc.bytes)-1] | 1
	}

	hc.lastBitsCount++
	return hc
}

func (hc HuffmanCode) flush() HuffmanCode {
	if hc.lastBitsCount == 0 {
		return hc
	}

	hc.bytes[len(hc.bytes)-1] = hc.bytes[len(hc.bytes)-1] << (8 - hc.lastBitsCount)
	return hc
}

type Node struct {
	Letter      byte
	Weight      int
	Left, Right *Node
}

func (n Node) isLeaf() bool {
	return n.Left == nil && n.Right == nil
}

func freqsToNodes(freqs map[byte]int) []Node {
	nodes := make([]Node, 0, len(freqs))

	for letter, freq := range freqs {
		nodes = append(nodes, Node{Letter: letter, Weight: freq})
	}

	return nodes
}

func BuildHuffmanTree(heap []Node) Node {
	if len(heap) == 1 {
		return heap[0]
	}

	// Descending order
	slices.SortFunc(heap, func(a Node, b Node) int {
		if b.Weight != a.Weight {
			return b.Weight - a.Weight
		}
		return int(a.Letter) - int(b.Letter)
	})

	var newIntNode Node

	for len(heap) > 1 {
		min1, min2 := heap[len(heap)-1], heap[len(heap)-2]
		newIntNode = Node{255, min1.Weight + min2.Weight, &min1, &min2}
		heap = heap[:len(heap)-1]

		for i := len(heap) - 2; i >= 0; i-- {
			if newIntNode.Weight <= heap[i].Weight {
				heap[i+1] = newIntNode
				break
			} else {
				heap[i+1] = heap[i]
				if i == 0 {
					heap[i] = newIntNode
				}
			}
		}
	}

	return newIntNode
}

func TraversePreorder(root *Node, codes map[byte]HuffmanCode, currentCode HuffmanCode) {
	if root == nil {
		return
	}

	if root.isLeaf() {
		if len(currentCode.bytes) == 0 {
			currentCode = currentCode.addBit(0)
		}
		currentCode.flush()
		codes[root.Letter] = currentCode
		return
	}

	TraversePreorder(root.Left, codes, currentCode.addBit(0))
	TraversePreorder(root.Right, codes, currentCode.addBit(1))
}
