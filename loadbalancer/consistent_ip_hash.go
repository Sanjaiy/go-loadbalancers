package loadbalancer

import (
	"fmt"
	"log"
	"slices"
)

type ConsistentIpHashBalancer struct {
	hashRing []uint32
	hashMap map[uint32]string
	nodes int
	servers map[string]bool
}

func NewConsistentIpHashBalancer(servers []string, node int) *ConsistentIpHashBalancer {
	cib := &ConsistentIpHashBalancer{
		servers: make(map[string]bool),
		hashMap: make(map[uint32]string),
		nodes: node,
	}
	
	for _, server := range servers {
		cib.AddServer(server)
	}

	return cib
}


func (cib *ConsistentIpHashBalancer) AddServer(server string) {
	if cib.servers[server] {
		log.Printf("Server %s already exist", server)
		return
	}

	cib.servers[server] = true

	for i := 1; i <= cib.nodes; i++ {
		hash := hashString(fmt.Sprintf("%s#%d", server, i))
		cib.hashRing = append(cib.hashRing, hash)
		cib.hashMap[hash] = server
	}

	slices.Sort(cib.hashRing)
}

func (cib *ConsistentIpHashBalancer) RemoveServer(server string) {
	if !cib.servers[server] {
		log.Printf("Server %s is not present", server)
		return
	}

	delete(cib.servers, server)

	newHashRing := []uint32{}
	for _, hash := range cib.hashRing {
		if cib.hashMap[hash] != server {
			newHashRing = append(newHashRing, hash)
		} else {
			delete(cib.hashMap, hash)
		}
	}

	cib.hashRing = newHashRing
}

func (cib *ConsistentIpHashBalancer) NextServer(clientIP string) string {
	if len(cib.hashRing) == 0 {
		return ""
	}

	hash := hashString(clientIP)

	for _, hashNode := range cib.hashRing {
		if hashNode >= hash {
			return cib.hashMap[hashNode]
		}
	}

	return cib.hashMap[cib.hashRing[0]]
}


func (cib ConsistentIpHashBalancer) ShowRing() {
	for _, ring := range cib.hashRing {
		fmt.Printf("%d -> %s\n", ring, cib.hashMap[ring])
	}
}