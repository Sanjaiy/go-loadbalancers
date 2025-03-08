package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

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

	// Least Connection Balancer
	balancer2 := loadbalancer.NewLeastConnectionBalancer(servers, 5)
	fmt.Printf("-----------------%s-----------------\n", "Least Connection")
	for i := range 20 {
		server := balancer2.NextServer()
		fmt.Printf("Request %d directed to %s\n", i+1, server)

		if i%2 == 0{
			for _, server := range servers {
				balancer2.ReleaseConnection(server)
			}
		}
	}
	fmt.Printf("-----------------%s-----------------\n\n", "Least Connection")

	// Least Response Time Balancer
	balancer3 := loadbalancer.NewLeastResponseTimeBalancer(servers, 5)
	for _, server := range balancer3.Servers {
		balancer3.PresetResponseTime(server.URL)
	}
	var wg sync.WaitGroup

	fmt.Printf("-----------------%s-----------------\n", "Least Response Time")
    for i := 0; i < 15; i++ {
		wg.Add(1)
		go func () {
			defer wg.Done()
			server, startTime := balancer3.NextServer()
			
			var responseDelay time.Duration
			switch server {
			case "server1.com":
				responseDelay = time.Duration(rand.Intn(35) + 30) * time.Millisecond
			case "server2.com":
				responseDelay = time.Duration(rand.Intn(35) + 30)  * time.Millisecond
			case "server3.com":
				responseDelay = time.Duration(rand.Intn(35) + 30) * time.Millisecond
			}
			
			endTime := startTime.Add(-responseDelay)
			balancer3.ReportCompletedRequest(server, endTime)
			fmt.Printf("Request %d used %s with simulated response time of %v\n", i+1, server, responseDelay)
		}()
    }
	wg.Wait()
	fmt.Printf("-----------------%s-----------------\n\n", "Least Response Time")
}