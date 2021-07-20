package etcd

import (
	"context"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdDiscovery struct {
	client *clientv3.Client
	lease  clientv3.Lease
	kv     clientv3.KV

	retryNum  int
	namespace string
	ctx       context.Context
	ttl       time.Duration
	stopChan  chan bool // for Deregister
}

func New(etcdAddr string, namespace string) (e *EtcdDiscovery, err error) {
	conf := clientv3.Config{}
	conf.EndPoint = strings.Split(etcdAddr, ",")
	client, err := clientv3.New(conf)
	if err != nil {
		return nil, err
	}
	e = &EtcdDiscovery{
		client:    client,
		kv:        clientv3.NewKV(client),
		retryNum:  5,
		namespace: namespace,
		ctx:       context.Background(),
		ttl:       time.Second * 10,
	}
	return e, nil
}

func (e *EtcdDiscovery) Register(ctx context.Context, path string, data string) (err error) {
	path = e.namespace + "/" + path
	if e.lease != nil {
		e.lease.Close()
	}
	e.lease = clientv3.NewLease(e.client)
	grant, err := e.lease.Grant(ctx, int64(e.ttl))
	if err != nil {
		return
	}
	_, err = e.client.Put(ctx, path, data, clientv3.WithLease(grant.ID))
	if err != nil {
		return
	}
	go e.heartBeat(ctx, grant.ID, path, data)
	return nil
	// leaseId, err := e.registerWithKV(ctx, path)
	// if r.lease != nil {
	// 	r.lease.Close()
	// }
	// r.lease = clientv3.NewLease(r.client)
	// leaseID, err := r.registerWithKV(ctx, key, value)
	// if err != nil {
	// 	return err
	// }

	// go r.heartBeat(ctx, leaseID, key, value)
	// return nil
}

// func (r *EtcdDiscovery) Deregister(ctx context.Context, service []byte]) error {
// 	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, service.Name, service.ID)
// 	value, err := marshal(service)
// 	if err != nil {
// 		return err
// 	}
// 	if r.lease != nil {
// 		r.lease.Close()
// 	}
// 	r.lease = clientv3.NewLease(r.client)
// 	leaseID, err := r.registerWithKV(ctx, key, value)
// 	if err != nil {
// 		return err
// 	}

// 	go r.heartBeat(ctx, leaseID, key, value)
// 	return nil
// }

// func (r *EtcdDiscovery) Watch(ctx context.Context, service []byte]) error {
// 	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, service.Name, service.ID)
// 	value, err := marshal(service)
// 	if err != nil {
// 		return err
// 	}
// 	if r.lease != nil {
// 		r.lease.Close()
// 	}
// 	r.lease = clientv3.NewLease(r.client)
// 	leaseID, err := r.registerWithKV(ctx, key, value)
// 	if err != nil {
// 		return err
// 	}

// 	go r.heartBeat(ctx, leaseID, key, value)
// 	return nil
// }

func (e *EtcdDiscovery) heartBeat(ctx context.Context, leaseId clientv3.LeaseID, key, val string) error {
	kac, err := e.client.KeepAlive(ctx, leaseId)
	if err != nil {
		leaseId = 0
	}
	for {
		if leaseId == 0 {

		}
	}
}

// 	go r.heartBeat(ctx, leaseID, key, value)
// 	return nil
// }
