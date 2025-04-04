// Name:Akeylah Hernandez
// Filename:main.go
// Assignment: Test #1
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Struct that defines each member in the scanner
type Definitions struct {
	Host    string `json:"host"`              //host of target
	Port    int    `json:"port"`              //the port number
	Check   bool   `json:"check"`             //check to see if the port is open or not
	Service string `json:"service,omitempty"` //banner grabbing result
}

// Struct that Display the scanned result
// Stores details about the scanning report
type DISPLAY struct {
	Targets    []string      `json:"targets"`
	TotalPorts int           `json:"total_ports"`
	Open       int           `json:"open"`
	DurationT  string        `json:"durationt"`
	Timeout    time.Duration `json:"timeout"`
	Workers    int           `json:"workers"`
	Range      string        `json:"range,omitempty"`
	Ports      []int         `json:"ports,omitempty"`
}

// Prints the scanner report and displays it
func Printer(printit DISPLAY) {
	fmt.Printf("Targets: %v \n", printit.Targets)
	if len(printit.Ports) > 0 {
		fmt.Printf("Ports: %v \n", printit.Ports)
	} else {
		fmt.Printf("Range: %s \n", printit.Range)
	}
	fmt.Printf("Ports that were Scanned: %d\n", printit.TotalPorts)
	fmt.Printf("Ports that are open: %d\n", printit.Open)
	fmt.Printf("Worker count: %d\n", printit.Workers)
	fmt.Printf("Timeout period: %s\n", printit.Timeout)
	fmt.Printf("Duration: %s\n", printit.DurationT)
}

// This function scans a list of targets for open ports
// Uses a WaitGroup and channels for concurrency
func worker(wg *sync.WaitGroup, tasks chan string, result chan Definitions, dialer net.Dialer, check bool) {
	defer wg.Done()
	for addr := range tasks {
		division := strings.Split(addr, ":") ///Gets the host and port
		if len(division) != 2 {
			continue
		}
		host, portStr := division[0], division[1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		conn, err := dialer.Dial("tcp", addr) //conneting to the adress of the port
		if err != nil {
			continue
		}

		resultDef := Definitions{ //Record an output for an open port
			Host:  host,
			Port:  port,
			Check: true,
		}

		if check { //banner grabbing implementation
			banner, err := GrabberHelper(conn, 1024, 2*time.Second)
			if err == nil {
				resultDef.Service = strings.TrimSpace(banner)
			}
		}
		conn.Close()
		result <- resultDef //Sends the result through the channel
	}
}

// Helper function that extract banners from ports that are open
func GrabberHelper(conn net.Conn, bufferSize int, timeout time.Duration) (string, error) {
	if conn == nil {
		return "",
			fmt.Errorf("no connection found") //error message if ran but no connection is established
	}
	err := conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return "", err
	}

	banner := make([]byte, bufferSize)
	Vari, err := conn.Read(banner)
	if err != nil {
		return "", err
	}
	return string(banner[:Vari]), nil
}

// Displays the progress of each port thats being scanned
func Progress(progress chan int, counter int) {
	first := time.Now()
	track := time.Now()
	var temp int
	for range progress {
		temp++
		now := time.Now()
		if now.Sub(track) > 500*time.Millisecond || temp == counter {
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

// Uses a slice to separate the ports
func ParseS(Flag string) ([]int, error) {
	if Flag == "" {
		return nil, nil
	}
	temp := strings.Split(Flag, ",")
	var ports []int
	for _, temp := range temp {
		port, err := strconv.Atoi(strings.TrimSpace(temp))
		if err != nil || port < 1 || port > 65535 {
			return nil, fmt.Errorf("invalid port: %s", temp)
		}
		ports = append(ports, port)
	}
	return ports, nil
}

// Main function that handles the whole scanning process
func main() {
	//Parse command line flags declarations
	startTime := time.Now()
	targets := flag.String("targets", "", "host")
	start := flag.Int("start-port", 1, "Port Number")
	end := flag.Int("end-port", 512, "Port Number")
	workers := flag.Int("worker", 5, "amount of workers")
	checkers := flag.Bool("booleancheck", false, "Banner grabbing")
	timeout := flag.Int("timeout", 5, "Timeout ")
	jsoutput := flag.Bool("json", false, "Enable jsoutput")
	specport := flag.String("ports", "", "List of points")
	flag.Parse()

	//Error message  if no target was specified
	fmt.Printf("Targets: %q\n", *targets)
	if *targets == "" {
		fmt.Println("No target specified! Try again")
		return
	}
	//  port scanning range
	var scannercheck []int
	if *specport != "" {
		ports, err := ParseS(*specport)
		if err != nil {
			fmt.Printf("Cannot parse ports: %v\n", err)
			return
		}
		scannercheck = ports
	} else {
		for port := *start; port <= *end; port++ {
			scannercheck = append(scannercheck, port)
		}
	}

	// Initialize synchronization and communication channels
	var wg sync.WaitGroup
	targetslist := strings.Split(*targets, ",")
	TotalPorts := len(targetslist) * len(scannercheck)
	tasks := make(chan string, 100)
	results := make(chan Definitions, *workers*2)
	progress := make(chan int, 100)
	go Progress(progress, TotalPorts)

	// Configure network dialer
	dialer := net.Dialer{
		Timeout: time.Duration(*timeout) * time.Second,
	}
	// goroutine to perform port scanning concurrently
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(&wg, tasks, results, dialer, *checkers)
	}

	go func() { //goes through the target host & list of ports, send task to worker and update the progress function
		for _, val := range targetslist {
			for _, port := range scannercheck {
				tasks <- net.JoinHostPort(val, strconv.Itoa(port))
				progress <- 1
			}
		}
		close(tasks)    //close channel for task sent
		close(progress) //close channel for progress
	}()
	wg.Wait()
	close(results) //close the channel

	// Collect results from workers
	var finals []Definitions
	for res := range results {
		finals = append(finals, res)
	}
	var Opens []Definitions
	var jsSummary []Definitions
	for _, result := range finals {
		//checks if the port is open
		if result.Check {
			Opens = append(Opens, result)
			if *checkers && result.Service != "" { //display the detected service

				fmt.Printf("%s:%d - IS OPEN (%s)\n", result.Host, result.Port, result.Service)
			} else {
				fmt.Printf("%s:%d - IS OPEN\n", result.Host, result.Port)
			}
		}
		if *jsoutput {
			jsSummary = append(jsSummary, result)
		}
	}
	//summary report of the scan
	Report := DISPLAY{
		Targets:    targetslist,
		TotalPorts: TotalPorts,
		Open:       len(Opens),
		DurationT:  time.Since(startTime).String(),
		Timeout:    time.Duration(*timeout) * time.Second,
		Workers:    *workers,
	}
	// Add scanned ports or port range to the report
	if *specport != "" {
		Report.Ports = scannercheck
	} else {
		Report.Range = fmt.Sprintf("%d-%d", *start, *end)
	}

	// output results in JSON format if requested
	if *jsoutput {
		print := struct {
			Results []Definitions `json:"results"`
			Summary DISPLAY       `json:"summary"`
		}{
			Results: jsSummary,
			Summary: Report,
		}

		// Convert results to JSON
		jsonData, err := json.MarshalIndent(print, "", "  ")
		if err != nil {
			fmt.Println("Error generating JSON output:", err)
			return
		}
		fmt.Println(string(jsonData))
	} else {
		Printer(Report) //print the scan report
	}
}
