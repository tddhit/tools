package msgqueue

import (
	"time"

	etcd "github.com/coreos/etcd/clientv3"

	diskqueue "github.com/tddhit/diskqueue/client"
	"github.com/tddhit/tools/msgqueue/option"
	"github.com/tddhit/wox/naming"
)

type DQConsumer struct {
	*diskqueue.Consumer
}

func NewDQConsumer(etcdClient *etcd.Client, opt option.DQConsumer, msgid string) (c *DQConsumer, err error) {
	c = &DQConsumer{}
	r := &naming.Resolver{
		Client:  etcdClient,
		Timeout: 2 * time.Second,
	}
	addrs := r.Resolve(opt.Registry)
	consumer := diskqueue.NewConsumer(opt.Topic, opt.Channel)
	for _, addr := range addrs {
		err = consumer.Connect(addr, msgid)
		if err != nil {
			return nil, err
		}
	}
	c.Consumer = consumer
	return c, nil
}
