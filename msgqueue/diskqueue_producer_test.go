package msgqueue

import (
	"strconv"
	"sync"
	"testing"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/tddhit/tools/log"
	"github.com/tddhit/tools/msgqueue/option"
)

func produce(i int) {
	cfg := etcd.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 2000 * time.Millisecond,
	}
	etcdClient, err := etcd.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	opt := option.DQProducer{
		Enable:        true,
		Registry:      "/nlpservice/diskqueue",
		RetryTimes:    3,
		RetryInterval: 5000,
	}
	p, err := NewDQProducer(etcdClient, opt)
	if err != nil {
		log.Fatal(err)
	}
	var j = 0
	for j = i * 10000; j < i*10000+10000; j++ {
		d := "hello" + strconv.Itoa(j)
		err := p.Publish("topic1", []byte(d))
		if err != nil {
			log.Fatal(err)
		}
	}
	p.Stop()
	log.Info(j)
}

func TestDQProducer(t *testing.T) {
	log.Init("producer.log", log.DEBUG)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			produce(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
