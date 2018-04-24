package msgqueue

import (
	"testing"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/tddhit/tools/log"
	"github.com/tddhit/tools/msgqueue/option"
)

func TestDQConsumer(t *testing.T) {
	log.Init("consumer.log", log.DEBUG)
	cfg := etcd.Config{
		Endpoints:   []string{"172.17.32.101:2379"},
		DialTimeout: 2000 * time.Millisecond,
	}
	log.Debug("!!")
	etcdClient, err := etcd.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("!!")
	opt := option.DQConsumer{
		Enable:   true,
		Registry: "/nlpservice/diskqueue",
		Topic:    "dict_candidate",
		Channel:  "wo",
	}
	log.Debug("!!")
	c, err := NewDQConsumer(etcdClient, opt, "0")
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("!!")
	for {
		msg := c.Pull()
		//time.Sleep(100 * time.Millisecond)
		log.Debug(string(msg))
	}
}