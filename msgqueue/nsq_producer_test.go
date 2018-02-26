package msgqueue

import (
	"testing"

	"github.com/tddhit/nlp/internal/conf"
)

func TestNsqProducer(t *testing.T) {
	conf := conf.Producer{
		Addrs:         []string{"localhost:4150"},
		RetryTimes:    3,
		RetryInterval: 5000,
		Topic:         make(map[string]bool),
	}
	conf.Topic["testTopic"] = true
	p := NewNsqProducer(conf)
	p.Publish("testTopic", []byte("testMessage1"))
	select {}
}
