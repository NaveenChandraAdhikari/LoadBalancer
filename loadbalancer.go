package loadbalancer

import (
	"context"
	"fmt"
	"load_balancer/algorithms"
	"load_balancer/health"
	"load_balancer/servers"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type LoadBalancer struct {

Servers []*servers.Server

Selector algorithms.Selector
HealthChecker  *health.HealthChecker
}



func NewLoadBalancer(
	servers []*servers.Server,
	selector algorithms.Selector,
) *LoadBalancer {
	return &LoadBalancer{
		Servers:  servers,
		Selector: selector,
		HealthChecker: health.NewHealthChecker(
			servers,
			2*time.Second,
		),
	}
}

func (lb *LoadBalancer ) handler(w http.ResponseWriter,r *http.Request){

	server :=lb.Selector.Next(lb.Servers)


if server == nil {
		http.Error(
			w,
			"No healthy servers available",
			http.StatusServiceUnavailable,
		)
		return
	}

	//Reverse Proxy → actually forwards the request to that server.
	//- Load Balancer → combines these responsibilities to distribute traffic.
	proxy := httputil.NewSingleHostReverseProxy(server.URL)

	proxy.ServeHTTP(w, r)

}


//Start() starts your HTTP load balancer on port 8080 and simultaneously starts the background health checker.
func (lb *LoadBalancer) Start(ctx context.Context) {

	//use my lb.handler func. Client → :8080 → lb.handler()
	http.HandleFunc("/", lb.handler)

	go lb.HealthChecker.Start(ctx)

	fmt.Println("Load balancer running on :8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("load balancer error:", err)
	}
}

//We use NewServer() to create and initialize a Server object from a URL.

func NewServer(rawURL string) *servers.Server {
	parsedURL, err := url.Parse(rawURL)

	if err != nil {
		panic(err)
	}

	server := &servers.Server{
		URL: parsedURL,
	}

	server.SetHealthy(true)

	return server
}