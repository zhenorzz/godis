package bptree

import (
	"errors"
	"fmt"
	"reflect"
)

var order = 4

type Tree struct {
	Root *Node
}

type Record struct {
	Value []byte
}

type Node struct {
	Pointers []interface{}
	Keys     []int
	Parent   *Node
	IsLeaf   bool
	NumKeys  int
	Next     *Node
}

func NewTree() *Tree {
	return &Tree{}
}

func (tree *Tree) Insert(key int, value []byte) error {
	var pointer *Record
	var leaf *Node

	if _, err := tree.Find(key, false); err == nil {
		return errors.New("key already exists")
	}

	pointer, err := makeRecord(value)
	if err != nil {
		return err
	}

	if tree.Root == nil {
		return tree.startNewTree(key, pointer)
	}

	leaf = tree.findLeaf(key, false)
	if leaf.NumKeys < order-1 {
		insertIntoLeaf(leaf, key, pointer)
		return nil
	}

	return tree.insertIntoLeafAfterSplitting(leaf, key, pointer)
}

func (tree *Tree) Find(key int, verbose bool) (*Record, error) {
	i := 0
	c := tree.findLeaf(key, verbose)
	if c == nil {
		return nil, errors.New("key not found")
	}
	for i = 0; i < c.NumKeys; i++ {
		if c.Keys[i] == key {
			break
		}
	}
	if i == c.NumKeys {
		return nil, errors.New("key not found")
	}
	r := c.Pointers[i].(*Record)
	return r, nil
}

func (tree *Tree) findLeaf(key int, verbose bool) *Node {
	i := 0
	c := tree.Root
	if c == nil {
		if verbose {
			fmt.Printf("Empty tree.\n")
		}
		return c
	}
	for !c.IsLeaf {
		if verbose {
			fmt.Printf("[")
			for i = 0; i < c.NumKeys-1; i++ {
				fmt.Printf("%d ", c.Keys[i])
			}
			fmt.Printf("%d]", c.Keys[i])
		}
		for i = 0; i < c.NumKeys; i++ {
			if key < c.Keys[i] {
				break
			}
		}
		if verbose {
			fmt.Printf("%d ->\n", i)
		}
		c = c.Pointers[i].(*Node)
	}
	if verbose {
		fmt.Printf("Leaf [")
		for i = 0; i < c.NumKeys-1; i++ {
			fmt.Printf("%d ", c.Keys[i])
		}
	}
	return c
}

func (tree *Tree) startNewTree(key int, pointer *Record) error {
	var err error
	tree.Root, err = makeLeaf()
	if err != nil {
		return err
	}
	tree.Root.Keys[0] = key
	tree.Root.Pointers[0] = pointer
	tree.Root.Pointers[order-1] = nil
	tree.Root.Parent = nil
	tree.Root.NumKeys += 1
	return nil
}

func makeLeaf() (*Node, error) {
	leaf, err := makeNode()
	if err != nil {
		return nil, err
	}
	leaf.IsLeaf = true
	return leaf, nil
}

func makeNode() (*Node, error) {
	newNode := new(Node)
	if newNode == nil {
		return nil, errors.New("error: Node creation")
	}
	newNode.Keys = make([]int, order-1)
	if newNode.Keys == nil {
		return nil, errors.New("error: New node keys array")
	}

	newNode.Pointers = make([]interface{}, order)
	if newNode.Keys == nil {
		return nil, errors.New("error: New node pointers array")
	}
	newNode.IsLeaf = false
	newNode.NumKeys = 0
	newNode.Parent = nil
	newNode.Next = nil
	return newNode, nil
}

func makeRecord(value []byte) (*Record, error) {
	newRecord := new(Record)
	if newRecord == nil {
		return nil, errors.New("error: record creation")
	} else {
		newRecord.Value = value
	}
	return newRecord, nil
}

func insertIntoLeaf(leaf *Node, key int, pointer *Record) {
	var i, insertionPoint int
	for insertionPoint < leaf.NumKeys && leaf.Keys[insertionPoint] < key {
		insertionPoint++
	}
	for i = leaf.NumKeys; i > insertionPoint; i-- {
		leaf.Keys[i] = leaf.Keys[i-1]
		leaf.Pointers[i] = leaf.Pointers[i-1]
	}
	leaf.Keys[insertionPoint] = key
	leaf.Pointers[insertionPoint] = pointer
	leaf.NumKeys++
	return
}

func (tree *Tree) insertIntoLeafAfterSplitting(leaf *Node, key int, pointer *Record) error {
	var newLeaf *Node
	var insertionIndex, split, newKey, i, j int
	var err error
	newLeaf, err = makeLeaf()
	if err != nil {
		return err
	}
	tempKeys := make([]int, order)
	if tempKeys == nil {
		return errors.New("error: Temporary keys array")
	}

	tempPointers := make([]interface{}, order)
	if tempPointers == nil {
		return errors.New("error: Temporary pointers array")
	}

	for insertionIndex < order-1 && leaf.Keys[insertionIndex] < key {
		insertionIndex++
	}

	for i = 0; i < leaf.NumKeys; i++ {
		if j == insertionIndex {
			j++
		}
		tempKeys[j] = leaf.Keys[i]
		tempPointers[j] = leaf.Pointers[i]
		j++
	}
	tempKeys[insertionIndex] = key
	tempPointers[insertionIndex] = pointer
	leaf.NumKeys = 0
	split = cut(order - 1)

	for i = 0; i < split; i++ {
		leaf.Pointers[i] = tempPointers[i]
		leaf.Keys[i] = tempKeys[i]
		leaf.NumKeys++
	}

	j = 0
	for i = split; i < order; i++ {
		newLeaf.Pointers[j] = tempPointers[i]
		newLeaf.Keys[j] = tempKeys[i]
		newLeaf.NumKeys++
		j++
	}

	newLeaf.Pointers[order-1] = leaf.Pointers[order-1]
	leaf.Pointers[order-1] = newLeaf

	for i = leaf.NumKeys; i < order-1; i++ {
		leaf.Pointers[i] = nil
	}

	for i = newLeaf.NumKeys; i < order-1; i++ {
		newLeaf.Pointers[i] = nil
	}

	newLeaf.Parent = leaf.Parent
	//child node include
	newKey = newLeaf.Keys[0]

	return tree.insertIntoParent(leaf, newKey, newLeaf)
}

func (tree *Tree) insertIntoParent(left *Node, key int, right *Node) error {
	var leftIndex int
	parent := left.Parent

	if parent == nil {
		return tree.insertIntoNewRoot(left, key, right)
	}

	leftIndex = getLeftIndex(parent, left)

	if parent.NumKeys < order-1 {
		insertIntoNode(parent, leftIndex, key, right)
		return nil
	}

	return tree.insertIntoNodeAfterSplitting(parent, leftIndex, key, right)

}

func getLeftIndex(parent, left *Node) int {
	leftIndex := 0
	for leftIndex <= parent.NumKeys && parent.Pointers[leftIndex] != left {
		leftIndex += 1
	}
	return leftIndex
}

func insertIntoNode(node *Node, leftIndex, key int, right *Node) {
	for i := node.NumKeys; i > leftIndex; i-- {
		node.Pointers[i+1] = node.Pointers[i]
		node.Keys[i] = node.Keys[i-1]
	}
	node.Pointers[leftIndex+1] = right
	node.Keys[leftIndex] = key
	node.NumKeys++
	return
}

func (tree *Tree) insertIntoNodeAfterSplitting(oldNode *Node, leftIndex, key int, right *Node) error {
	var i, j, split, kPrime int
	var newNode, child *Node
	var tempKeys []int
	var tempPointers []interface{}
	var err error

	tempPointers = make([]interface{}, order+1)
	if tempPointers == nil {
		return errors.New("error: Temporary point array for splitting nodes")
	}

	tempKeys = make([]int, order)
	if tempKeys == nil {
		return errors.New("error: Temporary keys array for splitting nodes")
	}

	for i = 0; i < oldNode.NumKeys+1; i++ {
		if j == leftIndex+1 {
			j++
		}
		tempPointers[j] = oldNode.Pointers[i]
		j++
	}

	j = 0
	for i = 0; i < oldNode.NumKeys; i++ {
		if j == leftIndex {
			j++
		}
		tempKeys[j] = oldNode.Keys[j]
		j++
	}

	tempPointers[leftIndex+1] = right
	tempKeys[leftIndex] = key

	split = cut(order)
	newNode, err = makeNode()
	if err != nil {
		return err
	}
	oldNode.NumKeys = 0
	for i = 0; i < split-1; i++ {
		oldNode.Pointers[i] = tempPointers[i]
		oldNode.Keys[i] = tempKeys[i]
		oldNode.NumKeys++
	}
	oldNode.Pointers[i] = tempPointers[i]
	kPrime = tempKeys[split-1]
	j = 0
	for i += 1; i < order; i++ {
		newNode.Pointers[j] = tempPointers
		newNode.Keys[j] = tempKeys[i]
		newNode.NumKeys++
		j++
	}
	newNode.Pointers[j] = tempPointers[i]
	newNode.Parent = oldNode.Parent
	for i = 0; i <= newNode.NumKeys; i++ {
		child = newNode.Pointers[i].(*Node)
		child.Parent = newNode
	}

	return tree.insertIntoParent(oldNode, kPrime, newNode)
}

func (tree *Tree) insertIntoNewRoot(left *Node, key int, right *Node) error {
	var err error
	tree.Root, err = makeNode()
	if err != nil {
		return err
	}
	tree.Root.Keys[0] = key
	tree.Root.Pointers[0] = left
	tree.Root.Pointers[1] = right
	tree.Root.NumKeys++
	tree.Root.Parent = nil
	left.Parent = tree.Root
	right.Parent = tree.Root
	return nil
}

func cut(length int) int {
	if length%2 == 0 {
		return length / 2
	}

	return length/2 + 1
}

func (tree *Tree) Delete(key int) error {
	keyRecord, err := tree.Find(key, false)
	if err != nil {
		return err
	}

	keyLeaf := tree.findLeaf(key, false)
	if keyRecord != nil && keyLeaf != nil {
		tree.deleteEntry(keyLeaf, key, keyRecord)
	}
	return nil
}
func (tree *Tree) deleteEntry(node *Node, key int, pointer interface{}) {
	var minKeys, neighbourIndex, kPrimeIndex, kPrime, capacity int
	var neighbour *Node

	node = removeEntryFromNode(node, key, pointer)
	if node == tree.Root {
		tree.adjustRoot()
		return
	}

	if node.IsLeaf {
		minKeys = cut(order - 1)
	} else {
		minKeys = cut(order) - 1
	}
	if node.NumKeys >= minKeys {
		return
	}
	neighbourIndex = getNeighbourIndex(node)

	if neighbourIndex == -1 {
		kPrimeIndex = 0
	} else {
		kPrimeIndex = neighbourIndex
	}

	kPrime = node.Parent.Keys[kPrimeIndex]

	if neighbourIndex == -1 {
		neighbour = node.Parent.Pointers[1].(*Node)
	} else {
		neighbour = node.Parent.Pointers[neighbourIndex].(*Node)
	}

	if node.IsLeaf {
		capacity = order
	} else {
		capacity = order - 1
	}

	if neighbour.NumKeys+node.NumKeys < capacity {
		tree.coalesceNodes(node, neighbour, neighbourIndex, kPrime)
		return
	} else {
		tree.redistributeNodes(node, neighbour, neighbourIndex, kPrimeIndex, kPrime)
	}
}

func removeEntryFromNode(node *Node, key int, pointer interface{}) *Node {
	var i, numPointers int

	for node.Keys[i] != key {
		i++
	}

	for i++; i < node.NumKeys; i++ {
		node.Keys[i-1] = node.Keys[i]
	}

	if node.IsLeaf {
		numPointers = node.NumKeys
	} else {
		numPointers = node.NumKeys + 1
	}

	i = 0
	for node.Pointers[i] != pointer {
		i++
	}
	for i++; i < numPointers; i++ {
		node.Pointers[i-1] = node.Pointers[i]
	}
	node.NumKeys--

	if node.IsLeaf {
		for i = node.NumKeys; i < order-1; i++ {
			node.Pointers[i] = nil
		}
	} else {
		for i = node.NumKeys + 1; i < order; i++ {
			node.Pointers[i] = nil
		}
	}
	return node
}

func (tree *Tree) adjustRoot() {
	var newRoot *Node
	if tree.Root.NumKeys > 0 {
		return
	}
	if !tree.Root.IsLeaf {
		newRoot = tree.Root.Pointers[0].(*Node)
		newRoot.Parent = nil
	} else {
		newRoot = nil
	}
	tree.Root = newRoot
	return
}

func getNeighbourIndex(node *Node) int {
	var i int
	for i = 0; i <= node.Parent.NumKeys; i++ {
		if reflect.DeepEqual(node.Parent.Pointers[i], node) {
			return i - 1
		}
	}
	return i
}

func (tree *Tree) coalesceNodes(node, neighbour *Node, neighbourIndex, kPrime int) {
	var i, j, neighbourInsertionIndex, nEnd int
	var tmp *Node

	if neighbourIndex == -1 {
		tmp = node
		node = neighbour
		neighbour = tmp
	}

	neighbourInsertionIndex = neighbour.NumKeys

	if !node.IsLeaf {
		neighbour.Keys[neighbourInsertionIndex] = kPrime
		neighbour.NumKeys++
		nEnd = node.NumKeys
		i = neighbourInsertionIndex + 1
		for j = 0; j < nEnd; j++ {
			neighbour.Keys[i] = node.Keys[j]
			neighbour.Pointers[i] = node.Pointers[j]
			neighbour.NumKeys++
			node.NumKeys--
			i++
		}
		neighbour.Pointers[i] = node.Pointers[j]

		for i = 0; i < neighbour.NumKeys+1; i++ {
			tmp = neighbour.Pointers[i].(*Node)
			tmp.Parent = neighbour
		}
	} else {
		i = neighbourInsertionIndex
		for j = 0; j < node.NumKeys; j++ {
			neighbour.Keys[i] = node.Keys[j]
			node.Pointers[i] = node.Pointers[j]
			neighbour.NumKeys++
		}
		neighbour.Pointers[order-1] = node.Pointers[order-1]
	}

	tree.deleteEntry(node.Parent, kPrime, node)
}

func (tree *Tree) redistributeNodes(node, neighbour *Node, neighbourIndex, kPrimeIndex, kPrime int) {
	var i int
	var tmp *Node

	if neighbourIndex != -1 {
		if !node.IsLeaf {
			node.Pointers[node.NumKeys+1] = node.Pointers[node.NumKeys]
		}
		for i = node.NumKeys; i > 0; i-- {
			node.Keys[i] = node.Keys[i-1]
			node.Pointers[i] = node.Pointers[i-1]
		}
		if !node.IsLeaf {
			node.Pointers[0] = neighbour.Pointers[neighbour.NumKeys]
			tmp = node.Pointers[0].(*Node)
			tmp.Parent = node
			neighbour.Pointers[neighbour.NumKeys] = nil
			node.Keys[0] = kPrime
			node.Parent.Keys[kPrimeIndex] = neighbour.Keys[neighbour.NumKeys-1]
		} else {
			node.Pointers[0] = neighbour.Pointers[neighbour.NumKeys-1]
			neighbour.Pointers[neighbour.NumKeys-1] = nil
			node.Keys[0] = neighbour.Keys[neighbour.NumKeys-1]
			node.Parent.Keys[kPrimeIndex] = node.Keys[0]
		}
	} else {
		if node.IsLeaf {
			node.Keys[node.NumKeys] = neighbour.Keys[0]
			node.Pointers[node.NumKeys] = neighbour.Pointers[0]
			node.Parent.Keys[kPrimeIndex] = neighbour.Keys[1]
		} else {
			node.Keys[node.NumKeys] = kPrime
			node.Pointers[node.NumKeys+1] = neighbour.Pointers[0]
			tmp, _ = node.Pointers[node.NumKeys+1].(*Node)
			tmp.Parent = node
			node.Parent.Keys[kPrimeIndex] = neighbour.Keys[0]
		}
		for i = 0; i < neighbour.NumKeys-1; i++ {
			neighbour.Keys[i] = neighbour.Keys[i+1]
			neighbour.Pointers[i] = neighbour.Pointers[i+1]
		}
		if !node.IsLeaf {
			neighbour.Pointers[i] = neighbour.Pointers[i+1]
		}
	}
	node.NumKeys += 1
	neighbour.NumKeys -= 1

	return
}
