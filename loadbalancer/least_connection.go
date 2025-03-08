package loadbalancer

import "sync"

type ServerState struct {
	URL string
	ActiveConnections int
	MaxConnections int
	mtx sync.Mutex
}

func NewServerState(url string, maxConnections int) *ServerState {
	return &ServerState{
		URL: url,
		ActiveConnections: 0,
		MaxConnections: maxConnections,
	}
}

func (ss *ServerState) IncrementConnections() {
	ss.mtx.Lock()
	defer ss.mtx.Unlock()
	ss.ActiveConnections++
}

func (ss *ServerState) DecrementConnections() {
	ss.mtx.Lock()
	defer ss.mtx.Unlock()
	ss.ActiveConnections--
}

func (ss *ServerState) GetConnections() int {
	ss.mtx.Lock()
	defer ss.mtx.Unlock()
	return ss.ActiveConnections
}

func (ss *ServerState) IsAvailable() bool {
	ss.mtx.Lock()
	defer ss.mtx.Unlock()
	return ss.ActiveConnections <= ss.MaxConnections
}

type LeastConnectionBalancer struct {
	servers []*ServerState
	mtx sync.Mutex
}

func NewLeastConnectionBalancer(serverURLs []string, maxConnections int) *LeastConnectionBalancer {
	servers := make([]*ServerState, len(serverURLs))
	for i, url := range serverURLs {
		servers[i] = NewServerState(url, maxConnections)
	}

	return &LeastConnectionBalancer{
		servers: servers,
	}
}

func (lcb *LeastConnectionBalancer) NextServer() string {
	lcb.mtx.Lock()
	defer lcb.mtx.Unlock()

	if len(lcb.servers) == 0 {
		return ""
	}

	var bestServer *ServerState
	minConnection := -1

	for _, server := range lcb.servers {
		if !server.IsAvailable() {
			continue
		}

		connections := server.GetConnections()
		if minConnection == -1 || connections < minConnection {
			minConnection = connections
			bestServer = server
		}
	}

	if bestServer == nil {
		return ""
	}

	bestServer.IncrementConnections()
	return bestServer.URL
}

func (lcb *LeastConnectionBalancer) ReleaseConnection(serverURL string) {
	lcb.mtx.Lock()
	defer lcb.mtx.Unlock()
	
	for _, server := range lcb.servers {
		if server.URL == serverURL {
			server.DecrementConnections()
			break
		}
	}
}