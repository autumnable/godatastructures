package avltree

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

type avlNode[K constraints.Ordered, V any] struct {
	count_, height_ int
	key             K
	left, right     *avlNode[K, V]
	value           V
}

func (this *avlNode[K, V]) balance() int {
	if this == nil {
		return 0
	}
	return this.left.height() - this.right.height()
}

func (this *avlNode[K, V]) ceiling(key K) (K, V, bool) {
	if this == nil {
		return key, zeroValue[V](), false
	} else if key > this.key {
		return this.right.ceiling(key)
	} else if key == this.key {
		return key, this.value, true
	} else if k, v, res := this.left.ceiling(key); res {
		return k, v, res
	} else {
		return this.key, this.value, true
	}
}

func (this *avlNode[K, V]) count() int {
	if this == nil {
		return 0
	}
	return this.count_
}

func (this *avlNode[K, V]) countGreater(key K) int {
	if this == nil {
		return 0
	} else if key <= this.key {
		return this.right.countGreater(key)
	}
	return 1 + this.right.count() + this.left.countGreater(key)
}

func (this *avlNode[K, V]) countLesser(key K) int {
	if this == nil {
		return 0
	} else if key >= this.key {
		return this.left.countGreater(key)
	}
	return 1 + this.left.count() + this.right.countLesser(key)
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

	this.count_--
	this.height_ = 1 + max(this.left.height(), this.right.height())
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

func (this *avlNode[K, V]) floor(key K) (K, V, bool) {
	if this == nil {
		return key, zeroValue[V](), false
	} else if key < this.key {
		return this.left.floor(key)
	} else if key == this.key {
		return key, this.value, true
	} else if k, v, res := this.right.floor(key); res {
		return k, v, res
	} else {
		return this.key, this.value, true
	}
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

func (this *avlNode[K, V]) has(value V) bool {
	return this != nil && (reflect.DeepEqual(this.value, value) || this.left.has(value) || this.right.has(value))
}

func (this *avlNode[K, V]) height() int {
	if this == nil {
		return 0
	}
	return this.height_
}

func (this *avlNode[K, V]) higher(key K) (K, V, bool) {
	if this == nil {
		return key, zeroValue[V](), false
	} else if key > this.key {
		return this.right.higher(key)
	} else if k, v, res := this.left.higher(key); res {
		return k, v, res
	} else {
		return this.key, this.value, true
	}
}

func (this *avlNode[K, V]) insert(key K, value V) (*avlNode[K, V], bool) {
	var res bool
	if this == nil {
		return &avlNode[K, V]{count_: 1, key: key, value: value}, true
	} else if key < this.key {
		this.left, res = this.left.insert(key, value)
	} else if key > this.key {
		this.right, res = this.right.insert(key, value)
	} else {
		this.value = value
		return this, false
	}

	this.count_++
	this.height_ = 1 + max(this.left.height(), this.right.height())

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

	this.count_ = 1 + this.left.count() + t2.count()
	y.count_ = 1 + this.count_ + y.right.count()
	y.left, this.right = this, t2

	this.height_ = max(this.left.height_, this.right.height_) + 1
	y.height_ = max(y.left.height_, y.right.height_) + 1
	return y
}

func (this *avlNode[K, V]) lower(key K) (K, V, bool) {
	if this == nil {
		return key, zeroValue[V](), false
	} else if key < this.key {
		return this.left.lower(key)
	} else if k, v, res := this.right.lower(key); res {
		return k, v, res
	} else {
		return this.key, this.value, true
	}
}

func (this *avlNode[K, V]) max() *avlNode[K, V] {
	for ; this.right != nil; this = this.right {
	}
	return this
}

func (this *avlNode[K, V]) min() *avlNode[K, V] {
	for ; this.left != nil; this = this.left {
	}
	return this
}

func (this *avlNode[K, V]) rightRotate() *avlNode[K, V] {
	x := this.left
	t2 := x.right

	this.count_ = 1 + this.right.count() + t2.count()
	x.count_ = 1 + this.count_ + x.left.count()
	x.right, this.left = this, t2

	this.height_ = max(this.left.height_, this.right.height_) + 1
	x.height_ = max(x.left.height_, x.right.height_) + 1
	return x
}

type AVLTree[K constraints.Ordered, V any] struct {
	root *avlNode[K, V]
}

func (this *AVLTree[K, V]) Add(key K, value V) bool {
	var res bool
	this.root, res = this.root.insert(key, value)
	return res
}

func (this *AVLTree[K, V]) AddAll(entries map[K]V) {
	for k, v := range entries {
		_ = this.Add(k, v)
	}
}

func (this *AVLTree[K, V]) Ceiling(key K) (K, V, bool) {
	return this.root.ceiling(key)
}

func (this *AVLTree[K, V]) Clear() {
	this.root = nil
}

func (this *AVLTree[K, V]) Count() int {
	return this.root.count()
}

func (this *AVLTree[K, V]) CountGreater(key K) int {
	return this.root.countGreater(key)
}

func (this *AVLTree[K, V]) CountGreaterOrEqual(key K) int {
	return this.root.count() - this.root.countLesser(key)
}

func (this *AVLTree[K, V]) CountLesser(key K) int {
	return this.root.countLesser(key)
}

func (this *AVLTree[K, V]) CountLesserOrEqual(key K) int {
	return this.root.count() - this.root.countGreater(key)
}

func (this *AVLTree[K, V]) First() (K, V, bool) {
	first := this.root.min()
	if first == nil {
		return zeroValue[K](), zeroValue[V](), false
	}
	return first.key, first.value, true
}

func (this *AVLTree[K, V]) Floor(key K) (K, V, bool) {
	return this.root.floor(key)
}

func (this *AVLTree[K, V]) Get(key K) (V, bool) {
	return this.root.get(key)
}

func (this *AVLTree[K, V]) Has(key K) bool {
	_, has := this.root.get(key)
	return has
}

func (this *AVLTree[K, V]) HasValue(value V) bool {
	return this.root.has(value)
}

func (this *AVLTree[K, V]) Height() int {
	return this.root.height()
}

func (this *AVLTree[K, V]) Higher(key K) (K, V, bool) {
	return this.root.higher(key)
}

func (this *AVLTree[K, V]) Keys() []K {
	nodes, res := this.preorder(), make([]K, 0, this.Count())
	for _, node := range nodes {
		res = append(res, node.key)
	}
	return res
}

func (this *AVLTree[K, V]) Last() (K, V, bool) {
	last := this.root.max()
	if last == nil {
		return zeroValue[K](), zeroValue[V](), false
	}
	return last.key, last.value, true
}

func (this *AVLTree[K, V]) Lower(key K) (K, V, bool) {
	return this.root.lower(key)
}

func (this *AVLTree[K, V]) PollFirst() (K, V, bool) {
	k, v, has := this.First()
	if !has {
		return k, v, has
	}
	_, _ = this.Remove(k)
	return k, v, true
}

func (this *AVLTree[K, V]) PollLast() (K, V, bool) {
	k, v, has := this.Last()
	if !has {
		return k, v, has
	}
	_, _ = this.Remove(k)
	return k, v, true
}

func (this *AVLTree[K, V]) Remove(key K) (V, bool) {
	if _, res := this.root.delete(key); res == nil {
		return zeroValue[V](), false
	} else {
		return res.value, true
	}
}

func (this *AVLTree[K, V]) ToMap() map[K]V {
	res := make(map[K]V, this.Count())
	for _, node := range this.preorder() {
		res[node.key] = node.value
	}
	return res
}

func (this *AVLTree[K, V]) Values() []V {
	res := make([]V, 0, this.Count())
	for _, node := range this.preorder() {
		res = append(res, node.value)
	}
	return res
}

func (this *AVLTree[K, V]) preorder() []*avlNode[K, V] {
	res := make([]*avlNode[K, V], 0, this.Count())
	var po func(*avlNode[K, V])
	po = func(node *avlNode[K, V]) {
		if node == nil {
			return
		}
		po(node.left)
		res = append(res, node)
		po(node.right)
	}
	return res
}

func zeroValue[T any]() T {
	var t T
	return t
}
