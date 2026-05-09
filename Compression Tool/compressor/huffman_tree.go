package compressor

import (
	"slices"
)

type HuffmanCode struct {
	Letter byte
	Code   string
}

type Node struct {
	Letter      byte
	Weight      int
	Left, Right *Node
}

func (n Node) isLeaf() bool {
	return n.Left == nil && n.Right == nil
}

func FreqsToNodes(freqs map[byte]int) []Node {
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

func TraversePreorder(root *Node, codes map[byte]string, currentCode string) {
	if root == nil {
		return
	}

	if root.isLeaf() {
		if currentCode == "" {
			currentCode = "0"
		}
		codes[root.Letter] = currentCode
		return
	}

	TraversePreorder(root.Left, codes, currentCode+"0")
	TraversePreorder(root.Right, codes, currentCode+"1")
}
