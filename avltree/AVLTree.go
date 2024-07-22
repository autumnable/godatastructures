package avltree

import (
	"golang.org/x/exp/constraints"
)

type avlNode[K constraints.Ordered, V any] struct {
	h           int
	key         K
	left, right *avlNode[K, V]
	value       V
}

func (node *avlNode[K, V]) balance() int {
	if node == nil {
		return 0
	}
	return node.left.height() - node.right.height()
}

func (node *avlNode[K, V]) delete(key K) (*avlNode[K, V], bool) {
	var res bool
	if node == nil {
		return nil, false
	} else if key < node.key {
		node.left, res = node.left.delete(key)
	} else if key > node.key {
		node.right, res = node.right.delete(key)
	} else if node.left == nil && node.right == nil {
		return nil, true
	} else if node.left == nil {
		node, res = node.right, true
	} else if node.right == nil {
		node, res = node.left, true
	} else {
		min := node.right.min()
		node.key, node.value = min.key, min.value
		node.right, res = node.right.delete(min.key)
	}

	node.h = 1 + max(node.left.height(), node.right.height())
	if b := node.balance(); b > 1 && node.left.balance() >= 0 {
		return rightRotate(node), res
	} else if b > 1 {
		node.left = leftRotate(node.left)
		return rightRotate(node), res
	} else if b < -1 && node.right.balance() <= 0 {
		return leftRotate(node), res
	} else if b < -1 {
		node.right = rightRotate(node.right)
		return leftRotate(node), res
	}
	return node, res
}

func (node *avlNode[K, V]) height() int {
	if node == nil {
		return 0
	}
	return node.h
}

func (node *avlNode[K, V]) insert(key K, value V) (*avlNode[K, V], bool) {
	var res bool
	if node == nil {
		return &avlNode[K, V]{key: key, value: value}, true
	} else if key < node.key {
		node.left, res = node.left.insert(key, value)
	} else if key > node.key {
		node.right, res = node.right.insert(key, value)
	} else {
		node.value = value
		return node, false
	}

	node.h = 1 + max(node.left.height(), node.right.height())

	if b := node.balance(); b > 1 && key < node.left.key {
		return rightRotate(node), res
	} else if b < -1 && key > node.right.key {
		return leftRotate(node), res
	} else if b > 1 && key > node.left.key {
		node.left = leftRotate(node.left)
		return rightRotate(node), res
	} else if b < -1 && key < node.right.key {
		node.right = rightRotate(node.right)
		return leftRotate(node), res
	}
	return node, res
}

func (node *avlNode[K, V]) min() *avlNode[K, V] {
	for node.left != nil {
		node = node.left
	}
	return node
}

type AVLTree[K constraints.Ordered, V any] struct {
	root *avlNode[K, V]
}

func (tree *AVLTree[K, V]) Height() int {
	return tree.root.height()
}

func leftRotate[K constraints.Ordered, V any](x *avlNode[K, V]) *avlNode[K, V] {
	var y *avlNode[K, V] = x.right
	var t2 *avlNode[K, V] = y.left

	// Perform rotation
	y.left = x
	x.right = t2

	// Update heights
	x.h = max(x.left.h, x.right.h) + 1
	y.h = max(y.left.h, y.right.h) + 1

	// Return new root
	return y
}

func rightRotate[K constraints.Ordered, V any](y *avlNode[K, V]) *avlNode[K, V] {
	var x *avlNode[K, V] = y.left
	var t2 *avlNode[K, V] = x.right

	// Perform rotation
	x.right = y
	y.left = t2

	// Update heights
	y.h = max(y.left.h, y.right.h) + 1
	x.h = max(x.left.h, x.right.h) + 1

	// Return new root
	return x
}
