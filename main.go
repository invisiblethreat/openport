package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/invisiblethreat/openport/sesame"
)

func main() {
	input := make(chan sesame.SingleTarget)
	output := make(chan sesame.Result)
	wg := sync.WaitGroup{}
	go sesame.Worker(input, output)
	go sesame.ResultHandler(output, &wg)

	addrs, err := sesame.ExpandAddrs(os.Args[1])
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	ports, err := sesame.ExpandPorts(os.Args[2])
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	protos, err := sesame.ExpandProtos(os.Args[3])
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	targets := sesame.AllTargets{Addrs: addrs, Ports: ports, Protos: protos}

	targets.Load(input, &wg)
	close(input)
	wg.Wait()
	close(output)
}
