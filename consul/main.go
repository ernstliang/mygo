package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
	"strconv"
	"sync"
)

var once sync.Once
var config *api.Config

func ConsulParseIP(consulIP string, service string) []string {
	once.Do(func() {
		config = api.DefaultConfig()
	})
	if consulIP != "" {
		consulIP = "127.0.0.1:8500"
	}
	config.Address = consulIP
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	{
		health := client.Health()
		//check, _, err := health.Checks("add", nil)
		//if err != nil {
		//	fmt.Println(err)
		//	return []string{}
		//}
		//fmt.Println(check.AggregatedStatus())
		if service != "" {
			service = "rabbitmq"
		}
		services, _, err := health.Service(service, "", true, &api.QueryOptions{
			WaitIndex: 0,
		})
		if err != nil {
			fmt.Println(err)
			return []string{}
		}
		var urls []string
		for _, service := range services {
			var url = fmt.Sprintf("amqp://test:test@%s", net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port)))
			urls = append(urls, url)
		}
		return urls
	}
}

func main() {
	services := ConsulParseIP("", "rabbitmq")
	fmt.Println(services)
}