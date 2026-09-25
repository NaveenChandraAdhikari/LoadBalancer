package algorithms

import "load_balancer/servers"

type Selector interface {
//Give me a list of servers, and I'll give you one server from that list,Then Next() chooses one server to handle the next request:.,an incoming HTTP request needs to be sent to one backend server, not all of them.
	Next(servers []*servers.Server) *servers.Server
}