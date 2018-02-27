package consistent

import (
	"crypto/sha1"
	"encoding/hex"
	"math"
	"sort"
	"strconv"
)

type HashRing struct {
	vprNum int // vnode num per rnode
	vnodes iVNodeArray
	rnodes map[string]RNode
}

type RNode struct {
	Id     string
	Weight int
}

type iVNode struct {
	key int
	id  string
}

type iVNodeArray []iVNode

func (n iVNodeArray) Len() int {
	return len(n)
}

func (n iVNodeArray) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n iVNodeArray) Less(i, j int) bool {
	return n[i].key < n[j].key
}

func (n iVNodeArray) Sort() {
	sort.Sort(n)
}

func NewHashRing(rnodes []RNode, vprNum int) *HashRing {
	hashRing := new(HashRing)
	hashRing.vprNum = vprNum
	hashRing.rnodes = make(map[string]RNode)
	for _, v := range rnodes {
		hashRing.rnodes[v.Id] = v
	}
	hashRing.adjust()
	return hashRing
}

func (h *HashRing) AddNode(rnode RNode) {
	h.rnodes[rnode.Id] = rnode
	h.adjust()
}

func (h *HashRing) RemoveNode(rnode RNode) {
	delete(h.rnodes, rnode.Id)
	h.adjust()
}

func (h *HashRing) GetNode(key string) RNode {
	k := hash(key)
	i := sort.Search(len(h.vnodes), func(i int) bool { return h.vnodes[i].key >= k })
	if i == len(h.vnodes) {
		i = 0
	}
	return h.rnodes[h.vnodes[i].id]
}

func (h *HashRing) adjust() {
	vnodesNum := len(h.rnodes) * h.vprNum
	totalWeight := 0
	for _, v := range h.rnodes {
		totalWeight += v.Weight
	}
	for k, v := range h.rnodes {
		j := int(math.Floor((float64(v.Weight) / float64(totalWeight)) * float64(vnodesNum)))
		for i := 0; i < j; i++ {
			key := hash(k + "#" + strconv.Itoa(i))
			h.vnodes = append(h.vnodes, iVNode{key, k})
		}
	}
	h.vnodes.Sort()
}

func hash(key string) int {
	rs := sha1.Sum([]byte(key))
	hexrs := hex.EncodeToString(rs[:])
	h, _ := strconv.ParseInt(hexrs[8:12], 16, 32)
	return int(h)
}
