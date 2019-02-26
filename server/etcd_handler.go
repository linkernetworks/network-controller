package main

import (
	"context"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type etcdCli struct {
	cli *clientv3.Client
}

func NewEtcdCli() *etcdCli {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.33.10:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &etcdCli{
		cli: cli,
	}
}

func (c server) put(ctx context.Context, key, value string) {
	if _, err := c.EtcdCli.cli.Put(ctx, key, value); err != nil {
		log.Printf("ETCD put data error beacause of %s\n", err)
	}
	log.Printf("ETCD put %s : %s n", key, value)
}
