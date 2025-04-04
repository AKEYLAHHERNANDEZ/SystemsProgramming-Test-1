# SystemsProgramming-Test-1
My Program: Its GO program that uses a TCP scanner, with flags and banner grabbing functionalities. 
Description Of Tools: The program takes a target host and a given range, and the ports that are to be scanned.
It created a few goroutines for the worker function, to check if there is any port open. If the banner grabber function is used, it attempts to get the service information. The results of the scan are displayed using the JSON format for the summary. There is also alot of conditions for error handling.

Instruction to build and run the program: 
Create an executable using the command: make build 
To run the scanner use the command: make run
To run a specific host & port use the command: ./portscanner -targets=  -ports=
Scan a range of ports: ./portscanner -targets= -start= -end=
To use banner grabbing use the command: ./portscanner -targets= -ports= -booleancheck
To display the JSON results use the command: ./portscanner -targets= -ports= -json

Or you can use the Makefile I created to run the program:
make build
make run
make clean
make test
make fmt


Sample output:

Target: "scanme.nmap.org"
Progress: 3/3 (100.0%) | Elapsed: 0s | Remaining: 0s
scanme.nmap.org:22 - IS OPEN (SSH)
scanme.nmap.org:80 - IS OPEN
Targets: ["scanme.nmap.org"]
Ports that were Scanned: 3
Ports that are open: 2
Worker count: 5
Timeout period: 5s
Duration: 2.3s




Homework #1 - Demo Video 
