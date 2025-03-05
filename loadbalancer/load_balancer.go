package loadbalancer

import "hash/fnv"

type LoadBalancer interface {
    NextServer(clientID string) string
}

func hashString(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}