package consistent

import (
	"testing"

	"github.com/tddhit/tools/log"
)

func TestHashRing(t *testing.T) {
	rnodes := make([]RNode, 0)
	rnodes = append(rnodes, RNode{"172.17.32.101:18900", 1})
	rnodes = append(rnodes, RNode{"172.17.32.101:18901", 2})
	rnodes = append(rnodes, RNode{"172.17.32.101:18902", 3})
	hashRing := NewHashRing(rnodes, 10)
	log.Debug(hashRing.GetNode("http://baike.baidu.com/item/爱情"))
	hashRing.AddNode(RNode{"127.0.0.6", 2})
	hashRing.RemoveNode(RNode{"127.0.0.5", 1})
}
