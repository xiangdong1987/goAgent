package test

import (
	"log"
	"nodeAgent/fun"
	"testing"
)

func TestEtcdGet(t *testing.T) {
	config := fun.EtcConfig{"localhost:2379", 5}
	c, err := fun.NewEtcClient(config)
	if err != nil {
		log.Print(err)
	}
	defer c.Client.Close()
	c.EtcGet("nodes/", true)
}

func TestEtcdPut(t *testing.T) {
	config := fun.EtcConfig{"localhost:2379", 5}
	c, err := fun.NewEtcClient(config)
	if err != nil {
		log.Print(err)
	}
	defer c.Client.Close()
	c.EtcPut("xdd", "123")
}

func TestEtcdWatch(t *testing.T) {
	config := fun.EtcConfig{"localhost:2379", 5}
	c, err := fun.NewEtcClient(config)
	if err != nil {
		log.Print(err)
	}
	defer c.Client.Close()
	c.EtcWatch("xdd")
}

func TestEtcdKeepAlive(t *testing.T) {
	config := fun.EtcConfig{"localhost:2379", 5}
	c, err := fun.NewEtcClient(config)
	if err != nil {
		log.Print(err)
	}
	defer c.Client.Close()
	ip, err := fun.ExternalIP()
	if err != nil {
		log.Println(err)
	}
	c.EtcKeepAlive("nodes/"+ip.String(), "online", 10)
}
func TestEtcdKeepAlive2(t *testing.T) {
	config := fun.EtcConfig{"localhost:2379", 5}
	c, err := fun.NewEtcClient(config)
	if err != nil {
		log.Print(err)
	}
	defer c.Client.Close()
	c.EtcKeepAlive("nodes/127.0.0.1", "online", 10)
}
func TestGetIp(t *testing.T) {
	ip, err := fun.ExternalIP()
	log.Println(ip, err)
}
