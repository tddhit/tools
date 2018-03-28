package msgqueue

import (
	"testing"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/tddhit/tools/log"
	"github.com/tddhit/tools/msgqueue/option"
)

func TestNsqProducer(t *testing.T) {
	cfg := etcd.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 2000 * time.Millisecond,
	}
	etcdClient, err := etcd.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	opt := option.NsqProducer{
		Enable:        true,
		Registry:      "/nlpservice/nsqd",
		RetryTimes:    3,
		RetryInterval: 5000,
	}
	p, err := NewNsqProducer(etcdClient, opt)
	if err != nil {
		log.Fatal(err)
	}
	p.Publish("testTopic", []byte("testMessage1"))
	select {}
}
