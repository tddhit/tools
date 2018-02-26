package msgqueue

import (
	"errors"
	"sync/atomic"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/nsqio/go-nsq"

	"github.com/tddhit/tools/log"
	"github.com/tddhit/tools/msgqueue/option"
	"github.com/tddhit/wox/naming"
)

var (
	errUnavailableProducer = errors.New("Unavailable Producer")
)

type NsqProducer struct {
	producers     []*nsq.Producer
	retryTimes    int
	retryInterval time.Duration
	counter       uint64
}

func NewNsqProducer(client *etcd.Client, opt option.NsqProducer) (p *NsqProducer, err error) {
	p = &NsqProducer{
		retryTimes:    opt.RetryTimes,
		retryInterval: time.Duration(opt.RetryInterval) * time.Millisecond,
	}
	cfg := nsq.NewConfig()
	r := &naming.Resolver{
		Client:  client,
		Timeout: 2 * time.Second,
	}
	addrs := r.Resolve(opt.Registry)
	for _, addr := range addrs {
		var producer *nsq.Producer
		if producer, err = nsq.NewProducer(addr, cfg); err != nil {
			log.Error(err)
			return
		}
		p.producers = append(p.producers, producer)
	}
	return
}

func (p *NsqProducer) Publish(topic string, body []byte) (err error) {
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
			if err = producer.DeferredPublish(topic, p.retryInterval, body); err != nil {
				continue
			}
			break
		}
		if err != nil {
			log.Errorf("type=msgqueue\tvendor=nsq\taddr=%s\ttopic=%s\tmsg=%s\tretry=%d\terr=%s\n",
				producer, topic, string(body), retryTimes, err)
			return
		}
	}
	log.Infof("type=msgqueue\tvendor=nsq\taddr=%s\ttopic=%s\tmsg=%s\tretry=%d\n",
		producer, topic, string(body), retryTimes)
	return
}
