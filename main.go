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
		if len(division) != 2 {
			continue
		}
		host, Variable := division[0], division[1]
		port, err := strconv.Atoi(Variable)
		if err != nil {
			continue
		}
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err == nil {
			result <- Definitions{
				Host:  host,
				Port:  port,
				Check: true,
			}
			defer conn.Close()
		}
		if check {
			banner, err := GrabberHelper(conn, 1024, 2*time.Second)
			if err != nil {
				banner = ""
			}
			result <- Definitions{
				Host:     host,
				Port:     port,
				Check:    true,
				Service:  banner,
				Protocol: protocol,
			}
		} else {
			result <- Definitions{
				Host:  host,
				Port:  port,
				Check: true,
			}
		}
	} 
} 
func GrabberHelper(conn net.Conn, bufferSize int, timeout time.Duration) (string, error) {
			banner := make([]byte,bufferSize)
			conn.SetReadDeadline(time.Now().Add(timeout))
			Vari, err := conn.Read(banner)
			if err != nil {
				return "",err
			}
			return string(banner[:Vari]), nil
			}


func main() {
	Flags := flag.String("Targets","")
	start := flag.Int("Start-port", 1, "Port Number")
	end := flag.Int("End-port", 512, "Port Number")
	workers := flag.Int("Start-point", 5, "Port Number")
	checkers := flag.Bool("Boolean check", false, "Banner grabbing")
	flag.Parse()
	
	var wg sync.WaitGroup
	tasks := make(chan string, 100)

    //target := "scanme.nmap.org"

	dialer := net.Dialer {
		Timeout: time.Duration(*timeout) *time.Second,
	}
	check = *checkers
	val := 100

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