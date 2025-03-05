package main

import (
	"fmt"

	"github.com/Sanjaiy/go-loadbalancer/loadbalancer"
)

func testWithClient(test string, balancer loadbalancer.LoadBalancer, servers, clients []string) {
	fmt.Printf("-----------------%s-----------------\n", test)
	for i, client := range clients {
		server := balancer.NextServer(client)
		fmt.Printf("Request %d from %s directed to %s\n", i+1, client, server)
	}
	fmt.Printf("-----------------%s-----------------\n\n", test)
}

func testWithOutClient(test string, balancer loadbalancer.LoadBalancer) {
	fmt.Printf("-----------------%s-----------------\n", test)
	for i := 0; i < 10; i++ {
		server := balancer.NextServer("")
		fmt.Printf("Request %d directed to %s\n", i+1, server)
	}
	fmt.Printf("-----------------%s-----------------\n\n", test)
} 

func main() {
	servers := []string {
		"server1.com",
		"server2.com",
		"server3.com",
	}
	clients := []string{"client1", "client2", "client3", "client1", "client2", "client4"}

	// Round Robin
	testWithOutClient("Round Robin", loadbalancer.NewRoundRobin(servers))

	// Sticky Round Robin
	testWithClient("Sticky Round Robin", loadbalancer.NewStickyRoundRobin(servers), servers, clients)

	// Weighted Round Robin
	testWithOutClient("Weighted Round Robin", loadbalancer.NewWeightedyRoundRobin(
		[]loadbalancer.Server{
			{URL: "server1.example.com", Weight: 2},
			{URL: "server2.example.com", Weight: 4},
			{URL: "server3.example.com", Weight: 6},
    	},
	))

	// IP Hash
	testWithClient("IP Hash", loadbalancer.NewIpHashBalancer(servers), servers, clients)

	// Consistent IP Hash
	balancer := loadbalancer.NewConsistentIpHashBalancer(servers, 5)
	balancer.ShowRing()
	testWithClient("Consistent IP Hash", balancer, servers, clients)
}