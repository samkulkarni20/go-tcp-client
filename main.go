package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"time"
)

func main() {
	threads := flag.Int("t", 100, "Number of threads")
	secondsDuration := flag.Int("d", 10*60, "Duration in seconds")
	port := flag.Int("p", 5777, "Port to connect to")

	// parse flags
	flag.Parse()

	// var wg sync.WaitGroup
	countCh := make(chan int, *threads)

	now := time.Now()
	endTime := now.Add(time.Duration(*secondsDuration) * time.Second)

	// // start threads
	for i := 1; i <= *threads; i++ {
		fmt.Printf("Starting thread %d\n", i)
		threadNum := i
		// wg.Add(1)
		go func() {
			// defer wg.Done()
			keepConnectingForDuration(countCh, *port, threadNum, endTime)
			// run(countCh, end, i)
		}()
	}
	// wg.Wait()

	var count int
	for i := 1; i <= *threads; i++ {
		count += <-countCh
	}

	fmt.Printf("Total connections: %d\n", count)
}

func run(countCh chan<- int, end time.Time, threadNum int) {
	fmt.Printf("Hello from thread %d \n", threadNum)
	count := 0
	now := time.Now()

	for now.Before(end) {
		time.Sleep(1 * time.Second)
		now = time.Now()
		count++
		fmt.Printf("Thread %d connected %d times\n", threadNum, count)
	}
	countCh <- count
}

func keepConnectingForDuration(countCh chan<- int, port, threadNum int, end time.Time) {
	count := 0
	now := time.Now()

	for now.Before(end) {
		connect(port, threadNum)
		now = time.Now()
		count++
	}
	countCh <- count
	fmt.Printf("Thread %d connected %d times\n", threadNum, count)
}

func connect(port, threadNum int) {
	strEcho := fmt.Sprintf("Hello from thread %d", threadNum)
	servAddr := fmt.Sprintf("localhost:%d", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		return
	}

	_, err = conn.Write([]byte(strEcho))
	if err != nil {
		println("Write to server failed:", err.Error())
		return
	}

	println("write to server = ", strEcho)

	reply := make([]byte, 50)

	n, err := conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		return
	}

	if bytes.NewBuffer(reply[:n]).String() != strEcho {
		println("Unexpected response from server. Received: %s", reply[:n])
	} else {
		println("reply from server=", string(reply[:n]))
	}
	println("Terminated")
	// time.Sleep(1 * time.Second)
	conn.Close()
}
