package main

import (
	"context"
	loadbalancer "load_balancer"

	"load_balancer/algorithms"
	"load_balancer/servers"
)

func main() {
	backendServers := []*servers.Server{
		loadbalancer.NewServer("http://localhost:8081"),
		loadbalancer.NewServer("http://localhost:8082"),
		loadbalancer.NewServer("http://localhost:8083"),
	}

	selector := algorithms.NewRoundRobin()

	//build the load balancer
	lb := loadbalancer.NewLoadBalancer(
		backendServers,
		selector,
	)

	ctx := context.Background()

	lb.Start(ctx)
}