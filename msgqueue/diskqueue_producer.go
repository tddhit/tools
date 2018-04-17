package msgqueue

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	etcd "github.com/coreos/etcd/clientv3"

	diskqueue "github.com/tddhit/diskqueue/client"
	"github.com/tddhit/tools/log"
	"github.com/tddhit/tools/msgqueue/option"
	"github.com/tddhit/wox/naming"
)

var (
	errUnavailableProducer = errors.New("Unavailable Producer")
)

type DQProducer struct {
	producers     []*diskqueue.Producer
	retryTimes    int
	retryInterval time.Duration
	counter       uint64
	wg            sync.WaitGroup
}

func NewDQProducer(client *etcd.Client, opt option.DQProducer) (p *DQProducer, err error) {
	p = &DQProducer{
		retryTimes:    opt.RetryTimes,
		retryInterval: time.Duration(opt.RetryInterval) * time.Millisecond,
	}
	r := &naming.Resolver{
		Client:  client,
		Timeout: 2 * time.Second,
	}
	addrs := r.Resolve(opt.Registry)
	log.Debug(addrs)
	for _, addr := range addrs {
		producer := diskqueue.NewProducer(addr)
		p.producers = append(p.producers, producer)
	}
	return
}

func (p *DQProducer) Publish(topic string, body []byte) (err error) {
	if len(p.producers) == 0 {
		err = errUnavailableProducer
		return
	}
	counter := atomic.AddUint64(&p.counter, 1)
	index := counter % uint64(len(p.producers))
	producer := p.producers[index]
	retryTimes := 0
	if err = producer.Publish(topic, body); err != nil {
		for retryTimes < p.retryTimes {
			retryTimes++
			if err = producer.Publish(topic, body); err != nil {
				continue
			}
			break
		}
		if err != nil {
			log.Errorf("type=msgqueue\tvendor=diskqueue\taddr=%s\ttopic=%s\tmsg=%s\tretry=%d\terr=%s\n",
				producer, topic, string(body), retryTimes, err)
			return
		}
	}
	log.Infof("type=msgqueue\tvendor=diskqueue\taddr=%s\ttopic=%s\tmsg=%s\tretry=%d\n",
		producer, topic, string(body), retryTimes)
	return
}

func (p *DQProducer) Stop() {
	for _, producer := range p.producers {
		p.wg.Add(1)
		go func() {
			producer.Stop()
			p.wg.Done()
		}()
	}
	p.wg.Wait()
}
