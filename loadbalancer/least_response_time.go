package loadbalancer

import (
	"math/rand"
	"sync"
	"time"
)

type ServerMetrics struct {
	URL string
	ActiveConnections int
	ResponseTimes []time.Duration
	MaxCapacity int
	TotalRequets int
	currentPointer int
	mtx sync.Mutex
}

func NewServerMetrics(url string, capacity int) *ServerMetrics {
	return &ServerMetrics{
		URL: url,
		ActiveConnections: 0,
		ResponseTimes: make([]time.Duration, 0),
		MaxCapacity: capacity,
		TotalRequets: 0,
		currentPointer: 0,
	}
}

func (sm *ServerMetrics) IncrementConnection() {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	sm.ActiveConnections++
	sm.TotalRequets++
}

func (sm *ServerMetrics) DecrementConnection() {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	if sm.ActiveConnections > 0 {
		sm.ActiveConnections--
	}
}

func (sm *ServerMetrics) RecordResponseTime(duration time.Duration) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	sm.ResponseTimes = append(sm.ResponseTimes, duration)
	
	if len(sm.ResponseTimes) > sm.MaxCapacity {
		sm.ResponseTimes = sm.ResponseTimes[1:]
	}
}

func (sm *ServerMetrics) GetAverageResponseTime() time.Duration {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	if len(sm.ResponseTimes) == 0 {
		return 0
	}

	var reponseTimeSum time.Duration

	for _, duration := range sm.ResponseTimes {
		reponseTimeSum += duration
	}

	return reponseTimeSum / time.Duration(len(sm.ResponseTimes))
}

func (sm *ServerMetrics) GetActiveConnections() int {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	return sm.ActiveConnections
}

type LeastResponseTimeBalancer struct {
	Servers []*ServerMetrics
	mtx sync.Mutex
}

func NewLeastResponseTimeBalancer(urls []string, sampleSize int) *LeastResponseTimeBalancer {
	var servers []*ServerMetrics
	for _, url := range urls {
		servers = append(servers, NewServerMetrics(url, sampleSize))
	}

	return &LeastResponseTimeBalancer{
		Servers: servers,
	}
}

func (lrb *LeastResponseTimeBalancer) NextServer() (string, time.Time) {
	lrb.mtx.Lock()
	defer lrb.mtx.Unlock()

	if len(lrb.Servers) == 0{
		return "", time.Now()
	}

	var bestServer *ServerMetrics
	var minResponseTime time.Duration = -1

	for _, server := range lrb.Servers {
		responseTime := server.GetAverageResponseTime()
		
		if minResponseTime == -1 || responseTime < minResponseTime {
			minResponseTime = responseTime
			bestServer = server
		}
	}

	if bestServer == nil {
		bestServer = lrb.Servers[0]
	}

	bestServer.IncrementConnection()

	return bestServer.URL, time.Now()
}

func (lrb *LeastResponseTimeBalancer) ReportCompletedRequest(url string, start time.Time) {
	duration := time.Since(start)

	lrb.mtx.Lock()
	defer lrb.mtx.Unlock()

	for _, server := range lrb.Servers {
		if server.URL == url {
			server.DecrementConnection()
			server.RecordResponseTime(duration)
			break
		}
	}
}

func (lrb *LeastResponseTimeBalancer) PresetResponseTime(serverURL string) {
    lrb.mtx.Lock()
    defer lrb.mtx.Unlock()
    
    for _, server := range lrb.Servers {
        if server.URL == serverURL {
            for i := 0; i < 5; i++ {
                server.RecordResponseTime(time.Duration(rand.Intn(35)+20)* time.Millisecond)
            }
        }
    }
}