//Name:Akeylah Hernandez
//Filename:main.go
//Assignemnt: Test #1
package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
	"encoding/json"
	"flag"
	"strings"
	"net"
)

type Definitions struct {
	Host  string `json:"host"`
	Port int	`json:"port"`
	Check bool	`json:"check"`
	Service string	`json:"service,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

func worker(wg *sync.WaitGroup, tasks chan string, result chan Definitions, dialer net.Dialer) {
	defer wg.Done()
    for addr := range tasks {
		division := strings.Split(addr, ":")
		host,Variable := division[0], division[1]
		port,_ := strconv.Atoi(Variable)

		conn, err := net.DialTimeout("tcp", addr,timeout)
		if err == nil {
			conn.Close()
			
			fmt.Printf("Connection to %s was successful\n", addr)
			banner := make([]byte,102)
		}
		backoff := time.Duration(1<<i) * time.Second
		fmt.Printf("Attempt %d to %s failed. Waiting %v...\n", i+1,  addr, backoff)
		time.Sleep(backoff)
	    }
		if !success {
			fmt.Printf("Failed to connect to %s after %d attempts\n", addr, maxRetries)
		}
	}
}

func main() {

	var wg sync.WaitGroup
	tasks := make(chan string, 100)

    target := "scanme.nmap.org"

	dialer := net.Dialer {
		Timeout: 5 * time.Second,
	}
  
	workers := 100

    for i := 1; i <= workers; i++ {
		wg.Add(1)
		go worker(&wg, tasks, dialer)
	}

	ports := 512

	for p := 1; p <= ports; p++ {
		port := strconv.Itoa(p)
        address := net.JoinHostPort(target, port)
		tasks <- address
	}
	close(tasks)
	wg.Wait()
}