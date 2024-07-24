package avltree

import "golang.org/x/exp/constraints"

type avlNode[K constraints.Ordered, V any] struct {
	h           int
	key         K
	left, right *avlNode[K, V]
	value       V
}

func (this *avlNode[K, V]) balance() int {
	if this == nil {
		return 0
	}
	return this.left.height() - this.right.height()
}

func (this *avlNode[K, V]) delete(key K) (res *avlNode[K, V], deleted *avlNode[K, V]) {
	if this == nil {
		return nil, nil
	} else if key < this.key {
		this.left, deleted = this.left.delete(key)
	} else if key > this.key {
		this.right, deleted = this.right.delete(key)
	} else if this.left == nil && this.right == nil {
		return nil, this
	} else if this.left == nil {
		deleted, this = this, this.right
	} else if this.right == nil {
		deleted, this = this, this.left
	} else {
		deleted = this.right.min()
		this.key, deleted.key = deleted.key, this.key
		this.value, deleted.value = deleted.value, this.value
		this.right, _ = this.right.delete(deleted.key)
	}

	this.h = 1 + max(this.left.height(), this.right.height())
	if b := this.balance(); b > 1 && this.left.balance() >= 0 {
		return this.rightRotate(), deleted
	} else if b > 1 {
		this.left = this.left.leftRotate()
		return this.rightRotate(), deleted
	} else if b < -1 && this.right.balance() <= 0 {
		return this.leftRotate(), deleted
	} else if b < -1 {
		this.right = this.right.rightRotate()
		return this.leftRotate(), deleted
	}
	return this, deleted
}

func (this *avlNode[K, V]) get(key K) (V, bool) {
	if this == nil {
		return zeroValue[V](), false
	} else if key < this.key {
		return this.left.get(key)
	} else if key > this.key {
		return this.right.get(key)
	}
	return this.value, true
}

func (this *avlNode[K, V]) height() int {
	if this == nil {
		return 0
	}
	return this.h
}

func (this *avlNode[K, V]) insert(key K, value V) (*avlNode[K, V], bool) {
	var res bool
	if this == nil {
		return &avlNode[K, V]{key: key, value: value}, true
	} else if key < this.key {
		this.left, res = this.left.insert(key, value)
	} else if key > this.key {
		this.right, res = this.right.insert(key, value)
	} else {
		this.value = value
		return this, false
	}

	this.h = 1 + max(this.left.height(), this.right.height())

	if b := this.balance(); b > 1 && key < this.left.key {
		return this.rightRotate(), res
	} else if b < -1 && key > this.right.key {
		return this.leftRotate(), res
	} else if b > 1 && key > this.left.key {
		this.left = this.left.leftRotate()
		return this.rightRotate(), res
	} else if b < -1 && key < this.right.key {
		this.right = this.right.rightRotate()
		return this.leftRotate(), res
	}
	return this, res
}

func (this *avlNode[K, V]) leftRotate() *avlNode[K, V] {
	y := this.right
	t2 := y.left
	y.left, this.right = this, t2

	this.h = max(this.left.h, this.right.h) + 1
	y.h = max(y.left.h, y.right.h) + 1
	return y
}

func (this *avlNode[K, V]) min() *avlNode[K, V] {
	for this.left != nil {
		this = this.left
	}
	return this
}

func (this *avlNode[K, V]) rightRotate() *avlNode[K, V] {
	x := this.left
	t2 := x.right
	x.right, this.left = this, t2

	this.h = max(this.left.h, this.right.h) + 1
	x.h = max(x.left.h, x.right.h) + 1
	return x
}

type AVLTree[K constraints.Ordered, V any] struct {
	count int
	root  *avlNode[K, V]
}

func (this *AVLTree[K, V]) Add(key K, value V) bool {
	var res bool
	this.root, res = this.root.insert(key, value)
	if res {
		this.count++
	}
	return res
}

func (this *AVLTree[K, V]) Count() int {
	return this.count
}

func (this *AVLTree[K, V]) Get(key K) (V, bool) {
	return this.root.get(key)
}

func (this *AVLTree[K, V]) Has(key K) bool {
	_, has := this.root.get(key)
	return has
}

func (this *AVLTree[K, V]) Height() int {
	return this.root.height()
}

func (this *AVLTree[K, V]) Remove(key K) (V, bool) {
	if _, res := this.root.delete(key); res == nil {
		return zeroValue[V](), false
	} else {
		this.count--
		return res.value, true
	}
}

func zeroValue[T any]() T {
	var t T
	return t
}
