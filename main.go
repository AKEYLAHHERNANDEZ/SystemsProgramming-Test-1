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

var (
	timeout time.Duration
	check bool
	service string
	protocol string
)

func worker(wg *sync.WaitGroup, tasks chan string, result chan Definitions, dialer net.Dialer) {
	defer wg.Done()
    for addr := range tasks {
		division := strings.Split(addr, ":")
		host,Variable := division[0], division[1]
		port,_ := strconv.Atoi(Variable)

		conn, err := net.DialTimeout("tcp", addr,timeout)
		if err == nil {
			conn.Close()
			
			const(
			buffersize = 1024
			timeout = 2 * time.Second
			)

			if check {
				banner,err := grabbanner (conn,buffersize,timeout)
				if err != nil {
					banner = ""
				}
				result <- Definitions{
					Host: host,
					Port: port,
					Check: true,
					Service: service,
					Protocol: protocol,
				}
			}
				else {
					result <- Definitions {
						Host: host,
						Port: port,
						Check: false,
					}	
			}

			func GrabberHelper(conn net.Coon, buffersizee int, timeout time.Duration) (string, error) {
			banner := make([byte,buffersize])
			CONN.SetReadDeadline(time.Now().Add(timeout))
			Vari, err :=conn.Read(bannerBuffer)
			if err != nil {
				return " ",err
			}
			return string(bannerBuffer[:vari],nil)
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