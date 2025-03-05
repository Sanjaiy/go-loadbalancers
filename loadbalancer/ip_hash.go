package loadbalancer

type IpHashBalancer struct {
	servers []string
}

func NewIpHashBalancer(servers []string) *IpHashBalancer {
	return &IpHashBalancer{
		servers: servers,
	}
}

func (ihb *IpHashBalancer) NextServer(clientIP string) string {
	if len(ihb.servers) == 0 {
		return ""
	}

	hash := hashString(clientIP)
	index := hash % uint32(len(ihb.servers))

	return ihb.servers[index]
}