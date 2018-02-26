package msgqueue

import (
	"testing"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/tools/msgqueue/option"
)

func do(body []byte) error {
	log.Debug(string(body))
	return nil
}

func TestNsqConsumer(t *testing.T) {
	conf := conf.Consumer{
		Id:     "testConsumer",
		Enable: true,
		Addrs:  []string{"localhost:4150"},
		Topic:  "testTopic",
	}
	c := NewNsqConsumer(conf)
	for msg := range c.Messages() {
		log.Debug(msg)
	}
}
