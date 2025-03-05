package loadbalancer

type Server struct {
	URL string
	Weight int
}

type WeightedRoundRobin struct {
	servers []Server
	currentWeight int
	currentIndex int
	maxWeight int
	gcd int
}

func NewWeightedyRoundRobin(servers []Server) *WeightedRoundRobin {
	maxWeight := 0
	var weights []int

	for _, s := range servers {
		if s.Weight > maxWeight {
			maxWeight = s.Weight
		}
		weights = append(weights, s.Weight)
	}
	
	return &WeightedRoundRobin{
		servers: servers,
		maxWeight: maxWeight,
		currentIndex: -1,
		currentWeight: 0,
		gcd: GCD(weights),
	}
}

func GCD(integers []int) int {
	if len(integers) == 0 {
        return 1
    }

	result := integers[0]
	for i:=1; i< len(integers);i++ {
		result = gcd(result, integers[i])
	}

	return result
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}

func (wrb *WeightedRoundRobin) NextServer(_ string) string {
	if len(wrb.servers) == 0 {
		return ""
	}

	for {
		wrb.currentIndex = (wrb.currentIndex + 1) % len(wrb.servers)
		if wrb.currentIndex == 0 {
			wrb.currentWeight -= wrb.gcd
			if wrb.currentWeight <= 0 {
				wrb.currentWeight = wrb.maxWeight
				if wrb.currentWeight == 0 {
					return ""
				}
			}
		}

		if wrb.servers[wrb.currentIndex].Weight >= wrb.currentWeight {
			return wrb.servers[wrb.currentIndex].URL
		}
	}
}