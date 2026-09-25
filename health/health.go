package health

import (
	"context"
	"fmt"
	"load_balancer/servers"
	"net/http"
	"time"
)


type HealthChecker struct {
	Servers []*servers.Server //all be servers
	Interval time.Duration //how often to check
	Client *http.Client //HTTP client used for /health
}

func NewHealthChecker (

	servers []*servers.Server,

    interval time.Duration,

) *HealthChecker {
	return &HealthChecker{
		Servers:  servers,
		Interval: interval,
		Client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}



func (hc *HealthChecker) check(server *servers.Server) {
	url := server.URL.String() + "/health"

	resp, err := hc.Client.Get(url)

	if err != nil {
		server.SetHealthy(false)
		fmt.Printf("%s unhealthy: %v\n", url, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		server.SetHealthy(true)
		fmt.Printf("%s healthy\n", url)
		return
	}

	server.SetHealthy(false)

	fmt.Printf(
		"%s unhealthy: status %d\n",
		url,
		resp.StatusCode,
	)
}

func (hc *HealthChecker) checkAll(){

	for _, server :=range hc.Servers{
		go hc.check(server)
	}


}


//Keep waiting forever. If the ticker fires, check all servers. If someone cancels the context, stop and exit

func (hc *HealthChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(hc.Interval) //ticket channel periodic work
	defer ticker.Stop() 

	for {
		select {

		case <-ticker.C:
			hc.checkAll()

		case <-ctx.Done():
			return

		}
	}
}