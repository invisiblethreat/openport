package sesame

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Worker checks the disposition of a host:port:protocol tuple and sends the
// result to the ResultHandler
func Worker(input chan SingleTarget, output chan Result) {
	for target := range input {
		conn, err := net.DialTimeout(target.Proto, target.toConn(), time.Duration(2*time.Second))
		if err != nil {
			output <- Result{Target: target, Open: false}
			continue
		}
		conn.Close()

		if err, ok := err.(*net.OpError); ok && err.Timeout() {
			fmt.Printf("Timeout error: %s\n", err)
			output <- Result{Target: target, Open: false}
			continue
		}

		if err != nil {
			// Log or report the error here
			fmt.Printf("Error: %s\n", err)
			output <- Result{Target: target, Open: false}
			continue
		}
		output <- Result{Target: target, Open: true}
	}
}

// ResultHandler collects scan results and processes them. It is also the
// blocking function to ensure that results are not orpahaned
func ResultHandler(input <-chan Result, wg *sync.WaitGroup) {
	for result := range input {
		fmt.Print(result.result())
		wg.Done()
	}
}
