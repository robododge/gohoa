package main

import (
	"fmt"
)

// Simple Trie playground to understand how it works and build a simple implementation
func main() {

	testTrie := InitTrie(NUM_TRIE)
	fmt.Println("---------- TRIE ----------")
	testTrie.Insert("4689")
	testTrie.Insert("4108")
	testTrie.Insert("4629")

	fmt.Printf("%+v\n", testTrie)

	fmt.Println("Search for 4689 = ", testTrie.Search("4689"))

	fmt.Println("Suggestions for 46 = ", testTrie.Suggestions("46"))
	fmt.Println("Suggestions for 4 = ", testTrie.Suggestions("4"))
	fmt.Println("Suggestions for 5 = ", testTrie.Suggestions("5"))

	fmt.Println("---------- TRIE strings----------")

	testTrie = InitTrie(TEXT_TRIE)
	testTrie.Insert("addrain")
	testTrie.Insert("bullings")
	testTrie.Insert("betterrup")
	testTrie.Insert("lancelot")
	testTrie.Insert("mackey")

	testTrie.Insert("maclace")
	testTrie.Insert("macbig")
	testTrie.Insert("muadib")
	testTrie.Insert("mallory")

	fmt.Println("Search for adrain = ", testTrie.Search("adrain"))
	fmt.Println("Search for b = ", testTrie.Suggestions("b"))
	fmt.Println("Search for m = ", testTrie.Suggestions("m"))
	fmt.Println("Search for mac = ", testTrie.Suggestions("mac"))

}

const (
	TEXT_TRIE = iota
	NUM_TRIE
)

type Node struct {
	Children []*Node
	isEnd    bool
}

func NewNode(tType int) *Node {
	if tType == NUM_TRIE {
		return &Node{Children: make([]*Node, 10)}
	}
	return &Node{Children: make([]*Node, 26)}

}

type Trie struct {
	root     *Node
	trieType int
	asciRune rune
}

func InitTrie(tType int) *Trie {
	asciiRune := '0'
	if tType == TEXT_TRIE {
		asciiRune = 'a'
	}

	return &Trie{NewNode(tType), tType, asciiRune}
}

func (t *Trie) Insert(key string) {
	node := t.root
	for _, c := range key {
		myIdx := c - t.asciRune
		if node.Children[myIdx] == nil {
			node.Children[myIdx] = NewNode(t.trieType)
		}
		node = node.Children[myIdx]
	}
	node.isEnd = true
}

func (t *Trie) Search(key string) bool {
	node := t.root
	for _, c := range key {
		myIdx := c - t.asciRune
		if node.Children[myIdx] == nil {
			return false
		}
		node = node.Children[myIdx]
	}
	return node.isEnd
}

func (t *Trie) Suggestions(key string) []string {
	node := t.root
	for _, c := range key {
		myIdx := c - t.asciRune
		if node.Children[myIdx] == nil {
			return nil
		}
		node = node.Children[myIdx]
	}
	return t.suggestionsFromNode(node, key)
}

func (t *Trie) suggestionsFromNode(node *Node, prefix string) []string {
	var suggestions []string
	if node.isEnd {
		suggestions = append(suggestions, prefix)
	}
	for i, child := range node.Children {
		if child != nil {
			childPrefix := prefix + string(t.asciRune+rune(i))
			childSuggestions := t.suggestionsFromNode(child, childPrefix)
			suggestions = append(suggestions, childSuggestions...)
		}
	}
	return suggestions
}
