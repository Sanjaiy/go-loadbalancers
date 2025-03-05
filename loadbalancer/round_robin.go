package loadbalancer

type RoundRobin struct {
	servers []string
	current int
}

func NewRoundRobin(servers []string) *RoundRobin {
	return &RoundRobin{
		servers: servers,
		current: 0,
	}
}

func (rb *RoundRobin) NextServer(_ string) string {
	if (len(rb.servers) == 0) {
		return ""
	}

	server := rb.servers[rb.current]
	rb.current = (rb.current + 1) % len(rb.servers)
	return server 
}
