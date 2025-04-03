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
	targets := flag.String("Targets","")
	start := flag.Int("Start-port", 1, "Port Number")
	end := flag.Int("End-port", 512, "Port Number")
	workers := flag.Int("Start-point", 5, "Port Number")
	checkers := flag.Bool("Boolean check", false, "Banner grabbing")
	timeout := flag.Int("timeout", 5, "Timeout ")
	jsoutput := flag.Bool("json", false, "Enable JSoutput")
	flag.Parse()
	
	if *targets == "" {
	fmt.Println("No targets specified")
	return
	}

	var wg sync.WaitGroup
	tasks := make(chan string, 100)

	dialer := net.Dialer {
		Timeout: time.Duration(*timeout) *time.Second,
	}
	check = *checkers

	results := make(chan Definitions, *workers) 
	targetslist := strings.Split(*targets, ",")

    for i := 0; i <= *workers; i++ {
		wg.Add(1)
		go worker(&wg, tasks, results, dialer)
	}

    go func(){
		for _, val := range targetslist {
			for port := *start; port <= *end; port++{
				tasks <- net.JoinHostPort(val,strconv.Itoa(port))
			}
		}
		close(tasks)
	}()

	go func(){
		wg.Wait()
		close(results)
	}()
var output []Definitions
for prints := range results {
		if prints.Check {
			if check {
				fmt.Printf("%s:%d - IS OPEN (%s)\n", prints.Host, prints.Port, prints.Service)
			} else {
				fmt.Printf("%s:%d - IS OPEN\n", prints.Host, prints.Port)
			}
		}
	}
}

if *jsoutput{
	jsonData,_:= json.MarshalIndent(output, "","")
	fmt.Println(string(jsonData))
}
