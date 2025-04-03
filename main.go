// Name:Akeylah Hernandez
// Filename:main.go
// Assignemnt: Test #1
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"text/scanner"
	"time"
)

type Definitions struct {
	Host  string `json:"host"`
	Port int	`json:"port"`
	Check bool	`json:"check"`
	Service string	`json:"service,omitempty"`
}

type DISPLAY struct{
	Targets      []string      `json:"targets"`
	TotalPorts   int           `json:"total_ports"`
	Open   int           		`json:"open"`
	DurationT  	 string        `json:"durationd"`
	Timeout      time.Duration `json:"timeout"`
	Workers      int           `json:"workers"`
	Range    string        	   `json:"range,omitempty"`
	Ports []int      		   `json:"ports,omitempty"`
}

func Printer(printit DISPLAY) {
	fmt.Println("Targets: %v \n", printit.Ports)
	if len(printit.Ports) > 0 {
		fmt.Println("Ports: %v \n", printit.Ports)
	}else {
		fmt.Println("Range: %s \n", printit.Range)
	}
	fmt.Println("Ports that were Scanned: %d\n", printit.TotalPorts)
	fmt.Println("Ports that are open: %d\n", printit.Open)
	fmt.Println("Worker count: %d\n", printit.Workers)
	fmt.Println("Timeout period: %s\n", printit.Timeout)
	fmt.Println("Duration: %s\n", printit.DurationT)	
}

func worker(wg *sync.WaitGroup, tasks chan string, result chan Definitions, dialer net.Dialer) {//fix
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
	conn, err := dialer.Dial("tcp", addr)
	if err == nil {
	continue
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
	if conn == nil {
		return "",
		fmt.Errorf("No connection Found")
	}
	err := conn.SetDeadline(time.Now().Add(timeout))
		if err != nil {
		return "",err
	}

	banner := make([]byte,bufferSize)
	Vari, err := conn.Read(banner)
	if err != nil {
	return "",err
	}
	return string(banner[:Vari]), nil
	}

func Progress(progress chan int, counter int){
	first := time.Now()
	track := time.Now()
	var temp int
	for range progress {
	temp++
	now := time.Now()
	if now.Sub(track) > 500 *time.Millisecond || temp == counter {
	amount := float64(temp) / float64(counter) * 100
	elapsed := time.Since(first)
	last := time.Duration(float64(elapsed)/float64(temp)*float64(counter-temp)) * time.Nanosecond
	fmt.Printf("\rProgress: %d/%d (%.1f%%) | Elapsed: %s | Remaining: %s",
	temp, counter, amount, elapsed.Round(time.Second), last.Round(time.Second))
	track = now
	if temp == counter {
		fmt.Println()
	}
}
}
}
func main() {
	targets := flag.String("Targets","","host")
	start := flag.Int("Start-port", 1, "Port Number")
	end := flag.Int("End-port", 512, "Port Number")
	workers := flag.Int("worker", 5, "amount of workers")
	checkers := flag.Bool("Boolean check", false, "Banner grabbing")
	timeout := flag.Int("timeout", 5, "Timeout ")
	jsoutput := flag.Bool("json", false, "Enable jsoutput")
	flag.Parse()
	
	if *targets == "" {
		fmt.Println("No target specified")
		return
	}
	
	var wg sync.WaitGroup
	tasks := make(chan string, 100)
	results := make(chan Definitions, *workers*2) 
	targetslist := strings.Split(*targets, ",")
	
	dialer := net.Dialer {
		Timeout: time.Duration(*timeout) * time.Second,
	}
	
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(&wg, tasks, results, dialer,*checkers)
	}
	
	go func(){
		for _, val := range targetslist {
		for port := *start; port <= *end; port++ {
		tasks <- net.JoinHostPort(val, strconv.Itoa(port))
		}
		}
		close(tasks)
	}()
	
	var Opens []Definitions
	var jsSummary []Definitions
	go func() {
		for prints := range results {
		if prints.Check {
			Opens = append(Opens, prints)
		if *checkers && prints.Service != "" {
		fmt.Printf("%s:%d - IS OPEN (%s)\n", prints.Host, prints.Port, prints.Service)
		}
		else {
		fmt.Printf("%s:%d - IS OPEN\n", prints.Host, prints.Port)
		}
	}
		if *jsoutput {
		output = append(output, prints)
		}
		}
	}()
	
	wg.Wait()
	close(results)
	
	Report := DISPLAY {
		Targets: targetslist,
		Port: len(targetslist) * (*end - *start +1),
		Open:  len(openPorts),
		Timeout: time.Duration(*timeout) * time.Second,
		Workers: *workers,
		Range: fmt.Sprintf("%d-%d", *start, *end),
	}
	if *jsoutput {
		jsonData, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(jsonData))
	}
}