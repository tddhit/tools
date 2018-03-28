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

func TestNsqConsumer(t *testing.T) {
	cfg := etcd.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 2000 * time.Millisecond,
	}
	etcdClient, err := etcd.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	opt := option.NsqConsumer{
		Enable:   true,
		Registry: "/nlpservice/nsqd",
		Topics:   []string{"testTopic"},
		Channel:  "wo",
	}
	c, err := NewNsqConsumer(etcdClient, opt)
	if err != nil {
		log.Fatal(err)
	}
	for msg := range c.Messages("testTopic") {
		log.Debug(string(msg.Body))
	}
}
