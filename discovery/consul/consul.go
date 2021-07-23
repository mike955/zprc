package consul

import "github.com/hashicorp/consul/api"

type ConsulDiscovery struct {
	client api.Client
}

func NewConsulDiscovery() (consul *ConsulDiscovery, err error) {
	api := api.NewClient()
	consul = &ConsulDiscovery{}
	return
}

func (c *ConsulDiscovery) Register() {

}

func (c *ConsulDiscovery) DeRegister() {

}

func (c *ConsulDiscovery) Watch() {

}
