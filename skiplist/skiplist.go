package skiplist

import (
	"bytes"
	"math/rand"
	"time"
	"unsafe"

	"github.com/tddhit/tools/log"
)

const MAXLEVEL = 32

type SkipList struct {
	level int
	size  int
	head  *Node
}

type Node struct {
	key     []byte
	value   []byte
	forward []*Node
}

func New() *SkipList {
	sk := &SkipList{
		head: &Node{
			forward: make([]*Node, MAXLEVEL),
		},
	}
	for i := 0; i < MAXLEVEL; i++ {
		sk.head.forward[i] = nil
	}
	rand.Seed(time.Now().UnixNano())
	return sk
}

func (sk *SkipList) Iterator() *SKIterator {
	ski := &SKIterator{
		sk:   sk,
		node: &Node{},
	}
	return ski
}

func (sk *SkipList) Show() {
	for i := sk.level - 1; i >= 0; i-- {
		p := sk.head.forward[i]
		log.Debug(i)
		for p != nil {
			log.Debugf("%s", string(p.key))
			p = p.forward[i]
		}
		log.Debug()
	}
}

func (sk *SkipList) Put(key, value []byte) {
	update := make([]*Node, MAXLEVEL)
	p := sk.head
	for i := sk.level - 1; i >= 0; i-- {
		for p.forward[i] != nil {
			r := bytes.Compare(key, p.forward[i].key)
			if r == -1 {
				break
			} else if r == 0 {
				sk.size = sk.size + len(value) - len(p.forward[i].value)
				p.forward[i].value = value
				return
			} else if r == 1 {
				p = p.forward[i]
			}
		}
		update[i] = p
	}
	level := randomLevel()
	node := &Node{
		key:     key,
		value:   value,
		forward: make([]*Node, level),
	}
	if level > sk.level {
		for i := sk.level; i < level; i++ {
			update[i] = sk.head
		}
		sk.level = level
	}
	for i := 0; i < level; i++ {
		node.forward[i] = update[i].forward[i]
		update[i].forward[i] = node
	}
	sk.size += len(key) + len(value) + int(unsafe.Sizeof(*node))
	sk.size += level * int(unsafe.Sizeof(node))
}

func (sk *SkipList) Get(key []byte) []byte {
	p := sk.head
	for i := sk.level - 1; i >= 0; i-- {
		for p.forward[i] != nil {
			r := bytes.Compare(key, p.forward[i].key)
			if r == -1 {
				break
			} else if r == 0 {
				return p.forward[i].value
			} else if r == 1 {
				p = p.forward[i]
			}
		}
	}
	return nil
}

func (sk *SkipList) Delete(key []byte) {
	update := make([]*Node, MAXLEVEL)
	var q *Node
	p := sk.head
	for i := sk.level - 1; i >= 0; i-- {
		for p.forward[i] != nil {
			r := bytes.Compare(key, p.forward[i].key)
			if r == -1 {
				break
			} else if r == 0 {
				q = p.forward[i]
				break
			} else if r == 1 {
				p = p.forward[i]
			}
		}
		update[i] = p
	}
	if q != nil {
		sk.size = sk.size - len(key) - len(p.value) - int(unsafe.Sizeof(*q))
		sk.size = sk.size - int(unsafe.Sizeof(q))*len(q.forward)
		for i := 0; i < len(q.forward); i++ {
			update[i].forward[i] = q.forward[i]
		}
	}
}

func (sk *SkipList) Size() int {
	return sk.size
}

func randomLevel() int {
	level := 1
	for i := 1; i < MAXLEVEL; i++ {
		if rand.Int()%2 == 0 {
			level++
		}
	}
	return level
}

type SKIterator struct {
	sk   *SkipList
	node *Node
}

func (ski *SKIterator) First() {
	if ski.sk != nil {
		ski.node = ski.sk.head.forward[0]
	}
}

func (ski *SKIterator) Next() {
	if ski.node.forward != nil {
		ski.node = ski.node.forward[0]
	}
}

func (ski *SKIterator) End() bool {
	if ski.node == nil {
		return true
	}
	return false
}

func (ski *SKIterator) Key() []byte {
	if ski.node != nil {
		return ski.node.key
	}
	return nil
}

func (ski *SKIterator) Value() []byte {
	if ski.node != nil {
		return ski.node.value
	}
	return nil
}
