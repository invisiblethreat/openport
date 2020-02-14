package sesame

import (
	"fmt"
	"sync"
)

// 2^16
const maxPort = 65536

// AllTargets holds the exploded arguments which are used for the Cartesian
// product to generate the set of atomic SingleTargets
type AllTargets struct {
	Addrs  []string
	Ports  []string
	Protos []string
}

// Load builds out atomic targets
func (at *AllTargets) Load(output chan SingleTarget, wg *sync.WaitGroup) {
	for _, proto := range at.Protos {
		for _, port := range at.Ports {
			for _, addr := range at.Addrs {
				output <- SingleTarget{Addr: addr, Port: port, Proto: proto}
				wg.Add(1)
			}
		}
	}
}

// SingleTarget is an atomic entity to attempt a connection
type SingleTarget struct {
	Addr  string `json:"addr"`
	Port  string `json:"port"`
	Proto string `json:"proto"`
}

// Result is the resulting disposition of contact with a host
type Result struct {
	Target SingleTarget `json:"host"`
	Open   bool         `json:"open"`
}

// toConn is a syntatic sugar function that is used in the dialer
func (t SingleTarget) toConn() string {
	return t.Addr + ":" + t.Port
}

// result is a an abstracted reporting function
func (r *Result) result() string {
	if r.Open {
		return fmt.Sprintf("%s is open for %s\n", r.Target.toConn(), r.Target.Proto)
	}
	return fmt.Sprintf("%s is closed for %s\n", r.Target.toConn(), r.Target.Proto)
}
