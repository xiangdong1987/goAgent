package fun

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"net"
	"time"
)

type EtcConfig struct {
	Address string
	TimeOut time.Duration
}
type EtcClient struct {
	Client *clientv3.Client
}

func NewEtcClient(config EtcConfig) (*EtcClient, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.Address},
		DialTimeout: config.TimeOut * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return nil, err
	}
	etcClient := &EtcClient{}
	etcClient.Client = cli
	return etcClient, nil
}

func (c *EtcClient) EtcPut(key string, value string) {

	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := c.Client.Put(ctx, key, value)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
}

func (c *EtcClient) EtcGet(key string, isPre bool) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	var err error
	var resp *clientv3.GetResponse
	if isPre {
		resp, err = c.Client.Get(ctx, key, clientv3.WithPrefix())
	} else {
		resp, err = c.Client.Get(ctx, key)
	}
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}

func (c *EtcClient) EtcWatch(key string) {
	rch := c.Client.Watch(context.Background(), key) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

func (c *EtcClient) EtcKeepAlive(key string, value string, leeseTime int64) {
	resp, err := c.Client.Grant(context.TODO(), leeseTime)
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.Client.Put(context.TODO(), key, value, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
	ch, kaerr := c.Client.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}
	for {
		ka := <-ch
		fmt.Println("ttl:", ka.TTL)
	}
}

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := GetIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func GetIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
