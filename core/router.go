package disel

import (
	"fmt"
	"strings"
)

var (
	NO_MISMATCH = -1
)

type Node struct {
	Value  interface{}
	Edges  map[byte]*Edge
	IsLeaf bool
}

func NewNode(isLeaf bool, value interface{}) *Node {
	if isLeaf {
		return &Node{
			Value:  value,
			Edges:  make(map[byte]*Edge),
			IsLeaf: isLeaf,
		}
	}
	return &Node{
		Edges:  make(map[byte]*Edge),
		IsLeaf: isLeaf,
	}
}

type Edge struct {
	Label string
	Next  *Node
}

func NewEdge(key string, next *Node, value interface{}) *Edge {
	return &Edge{
		Label: key,
		Next:  NewNode(true, value),
	}
}

func (n *Node) TotalEdges() int {
	return len(n.Edges)
}

func (n *Node) AddEdge(key string, next *Node, value interface{}) {
	n.Edges[key[0]] = NewEdge(key, next, value)
}

func (n *Node) GetTransition(transitionCharacter byte) *Edge {
	return n.Edges[transitionCharacter]
}

type RadixTree struct {
	Root *Node
}

func NewRadixTree() RadixTree {
	return RadixTree{
		Root: NewNode(false, nil),
	}
}

func (t *RadixTree) Insert(path string, handler interface{}) {
	current := t.Root
	currIdx := 0
	for {
		if currIdx >= len(path) {
			break
		}

		transitionChar := path[currIdx]
		currentEdge := current.GetTransition(transitionChar)
		currStr := path[currIdx:]

		if currentEdge == nil {
			current.Edges[transitionChar] = NewEdge(currStr, NewNode(true, handler), handler)
			break
		}

		splitIndex := GetFirstMismatchLetter(currStr, currentEdge.Label)
		if splitIndex == NO_MISMATCH {
			//The edge and leftover string are the same length
			//so finish and update the next node as a word node
			if len(currStr) == len(currentEdge.Label) {
				currentEdge.Next.IsLeaf = true
				break
			} else if len(currStr) < len(currentEdge.Label) {
				//The leftover word is a prefix to the edge string, so split
				suffix := currentEdge.Label[len(currStr):]
				currentEdge.Label = currStr
				newNext := NewNode(true, handler)
				afterNewNext := currentEdge.Next
				currentEdge.Next = newNext
				newNext.AddEdge(suffix, afterNewNext, afterNewNext.Value)
				break
			} else {
				splitIndex = len(currentEdge.Label)
			}
		} else {
			suffix := currentEdge.Label[splitIndex:]
			currentEdge.Label = currentEdge.Label[:splitIndex]
			prevNext := currentEdge.Next
			currentEdge.Next = NewNode(false, nil)
			currentEdge.Next.AddEdge(suffix, prevNext, prevNext.Value)
		}
		current = currentEdge.Next
		currIdx += splitIndex
	}
}

func (t *RadixTree) PrintAllWords() {
	printAllWords(t.Root, "")
}

func printAllWords(current *Node, result string) {
	if current.IsLeaf {
		fmt.Println(result)
	}

	for _, edge := range current.Edges {
		printAllWords(edge.Next, result+edge.Label)
	}
}

func (t *RadixTree) Search(path string) (*Node, bool) {
	current := t.Root
	currIndex := 0

	for {
		if currIndex >= len(path) {
			break
		}
		transitionChar := path[currIndex]
		edge := current.GetTransition(transitionChar)
		if edge == nil {
			return nil, false
		}
		currSubstring := path[currIndex:]
		if !strings.HasPrefix(currSubstring, edge.Label) {
			return nil, false
		}
		currIndex += len(edge.Label)
		current = edge.Next
	}
	if current.IsLeaf {
		return current, current.IsLeaf
	}
	return nil, current.IsLeaf
}

func GetFirstMismatchLetter(s1, s2 string) int {
	n1 := len(s1)
	n2 := len(s2)
	var minLen int

	if n1 > n2 {
		minLen = n2
	} else {
		minLen = n1
	}

	for i := 0; i < minLen; i++ {
		if s1[i] != s2[i] {
			return i
		}
	}
	return NO_MISMATCH
}
