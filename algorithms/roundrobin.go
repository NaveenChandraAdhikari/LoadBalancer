package algorithms

import (
	"load_balancer/servers"
	"sync"
)

type RoundRobin struct {

	current int
	mu sync.Mutex

}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}


//Give me the list of servers, and I will return the next healthy server according to the round-robin rotation

func (rr *RoundRobin) Next(servers []*servers.Server) *servers.Server {

	/*
Request A → 🔒 → choose Server 1 → current = 1 → 🔓
Request B → waits
              ↓
           🔒 → choose Server 2 → current = 2 → 🔓
	*/
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if len(servers) == 0 {
		return nil
	}

	for i := 0; i < len(servers); i++ {
		server := servers[rr.current]
		rr.current = (rr.current + 1) % len(servers)

		if server.IsHealthy() {
			return server
		}
	}

	return nil
}
