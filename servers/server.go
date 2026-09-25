package servers

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type Server struct {
	URL  *url.URL
	mu sync.RWMutex //protects shared Server state,lock that protects shared data when multiple goroutines access it
	isHealthy bool
}

func(s *Server) IsHealthy() bool{

	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isHealthy

}

func (s *Server)SetHealthy(healthy bool){
	s.mu.Lock()
	defer s.mu.Unlock()

	s.isHealthy=healthy
}


func Start(port string, id string) {
	//servemux is a router that decides which func should handle an incoming HTTP request.
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from %s\n", id)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	fmt.Printf("%s running on :%s\n", id, port)

	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		fmt.Println("server error:", err)
	}
}