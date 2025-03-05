package loadbalancer

type StickyRoundRobin struct {
	servers []string
	current int
	client map[string]string
}

func NewStickyRoundRobin(servers []string) *StickyRoundRobin {
	return &StickyRoundRobin{
		servers: servers,
		current: 0,
		client: make(map[string]string),
	}
}

func (srb *StickyRoundRobin) NextServer(client string) string {
	if len(srb.servers) == 0 {
		return ""
	}

	if server, exists := srb.client[client]; exists {
		for _, s := range srb.servers {
			if server == s {
				return server
			}
		}
	}

	server := srb.servers[srb.current]
	srb.current = (srb.current + 1) % len(srb.servers)
	
	srb.client[client] = server

	return server
}