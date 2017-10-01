package main

import (
    "fmt"
    "hello/ipc"
    "net"
    "os"
    "runtime"
	"bufio"
	"io"
    "bytes"
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

// Handles incoming requests.
func handleRequest(conn net.Conn) {
    rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
    defer conn.Close()
    
    var buffer bytes.Buffer
    readLoop:
    for {
        cmd, err := rw.ReadString('\n')
        buffer.WriteString(cmd)
		switch {
		case err == io.EOF:
			fmt.Println("Reached EOF - close this connection.\n   ---")
			break readLoop 
		case err != nil:
			fmt.Println("\nError reading command. Got: '"+cmd+"'\n", err)
            if err != nil {
                fmt.Println("Cannot write to connection.\n", err)
                
            }
			break readLoop 
        }
    }
    
    handleStrings(buffer.String(), rw)
}
func handleStrings(cmd string, rw *bufio.ReadWriter) {
    fmt.Println("handle: " + cmd)
    _, err := rw.WriteString("Thank you.\n")
    if err != nil {
        fmt.Println("Cannot write to connection.\n", err)
        return
    }
    err = rw.Flush()
    if err != nil {
        fmt.Println("Flush failed.", err)
    }
}