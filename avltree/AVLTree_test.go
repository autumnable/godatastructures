package avltree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	f float64
	i int
	s string
}

func TestAVLTree(t *testing.T) {
	is := assert.New(t)
	tree := AVLTree[float64, testStruct]{}

	// left rotate
	tree.Add(1, testStruct{f: 1.1})
	tree.Add(2, testStruct{i: 1})
	tree.Add(3, testStruct{s: "hi"})
	is.NotEqual(nil, tree.root.left)
	is.NotEqual(nil, tree.root.right)
	is.Equal(nil, tree.root.left.left)
	is.Equal(nil, tree.root.left.right)
	is.Equal(nil, tree.root.right.left)
	is.Equal(nil, tree.root.right.right)

	// right rotate
	tree.Clear()
	tree.Add(-1, testStruct{f: 1.1})
	tree.Add(-2, testStruct{i: 1})
	tree.Add(-3, testStruct{s: "hi"})
	is.NotEqual(nil, tree.root.left)
	is.NotEqual(nil, tree.root.right)
	is.Equal(nil, tree.root.left.left)
	is.Equal(nil, tree.root.left.right)
	is.Equal(nil, tree.root.right.left)
	is.Equal(nil, tree.root.right.right)
}
