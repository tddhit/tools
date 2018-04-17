package msgqueue

import (
	"testing"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/tddhit/tools/log"
	"github.com/tddhit/tools/msgqueue/option"
)

func do(body []byte) error {
	log.Debug(string(body))
	return nil
}

func TestDQConsumer(t *testing.T) {
	log.Init("consumer.log", log.DEBUG)
	cfg := etcd.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 2000 * time.Millisecond,
	}
	etcdClient, err := etcd.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	opt := option.DQConsumer{
		Enable:   true,
		Registry: "/nlpservice/diskqueue",
		Topic:    "topic1",
		Channel:  "wo2",
	}
	c, err := NewDQConsumer(etcdClient, opt, "0")
	if err != nil {
		log.Fatal(err)
	}
	for {
		msg := c.Pull()
		log.Debug(string(msg))
	}
}
