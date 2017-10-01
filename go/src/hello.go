package main

import (
	"bufio"
	"bytes"
	"fmt"
	"hello/ipc"
	"io"
	"net"
	"os"
	"runtime"
)

func main() {
	fmt.Printf("hello, world\n")
	// Listen on a random endpoint.
	endpoint := fmt.Sprintf("secrettunnel")
	if runtime.GOOS == "windows" {
		endpoint = `\\.\pipe\` + endpoint
	} else {
		endpoint = os.TempDir() + "/" + endpoint
	}
	listener, err := ipc.CreateIPCListener(endpoint)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on ipc endpoint: " + endpoint)
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleRequest(conn)
	}
}

func write2js(msgWriter *bufio.Writer, go2jsMessage <-chan string) {
	for {
		// I will run forever
		msg := <-go2jsMessage
		fmt.Println("write to js: " + msg)
		msgWriter.Write([]byte(msg))
		msgWriter.Flush()
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	msgReader := bufio.NewReader(conn)
	msgWriter := bufio.NewWriter(conn)
	//rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	go2jsMessage := make(chan string)
	_, err := msgWriter.WriteString("channel set up!")
	if err != nil {
		fmt.Println("Cannot write to connection.\n", err)
		return
	}
	msgWriter.Flush()
	go write2js(msgWriter, go2jsMessage)

	var buffer bytes.Buffer
readLoop:
	for {
		line, err := msgReader.ReadString('\n')
		if line == "\n" {
			// one piece of message received
			message := buffer.String()
			buffer.Reset()
			handleStrings(message, go2jsMessage)
		}
		buffer.WriteString(line)
		switch {
		case err == io.EOF:
			fmt.Println("Reached EOF - close this connection.")
			break readLoop
		case err != nil:
			fmt.Println("\nError reading command. Got: '"+line+"'\n", err)
			if err != nil {
				fmt.Println("Cannot write to connection.\n", err)
			}
			break readLoop
		}
	}

}
func handleStrings(cmd string, go2jsMessage chan<- string) {
	fmt.Println("handle: " + cmd)
	go2jsMessage <- "Thank you." + cmd
	// if err != nil {
	// 	fmt.Println("Cannot write to connection.\n", err)
	// 	return
	// }
	// err = rw.Flush()
	// if err != nil {
	// 	fmt.Println("Flush failed.", err)
	// }
}
