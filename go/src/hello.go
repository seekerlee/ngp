package main

import (
	"bufio"
	"fmt"
	"ipc"
	"io"
	"net"
	"os"
	"runtime"
	"strings"
	"encoding/binary"
)

func main() {
	// Listen on an endpoint.
	endpoint := fmt.Sprintf("secrettunnel")
	if runtime.GOOS == "windows" {
		endpoint = `\\.\pipe\` + endpoint
	} else {
		endpoint = strings.TrimSuffix(os.TempDir(), "/") + "/" + endpoint
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

func byteArr2Uint(p []byte) uint32 {
    data := binary.BigEndian.Uint32(p)
	fmt.Println(data)
	return data
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	msgReader := bufio.NewReader(conn)
	msgWriter := bufio.NewWriter(conn)
	defer conn.Close()

	sizebytes := make([]byte, 4)
	line, err := msgReader.Read(sizebytes)
	switch {
	case err == io.EOF:
		fmt.Println("Reached EOF - close this connection.")
		return
	case err != nil:
		//fmt.Println("\nError reading command. Got: '"+line+"'\n", err)
		if err != nil {
			fmt.Println("Cannot write to connection.\n", err)
		}
		return
	}
	fmt.Println(line)
	bodyLen := byteArr2Uint(sizebytes)
	bodyBytes := make([]byte, bodyLen)
	_, err2 := msgReader.Read(bodyBytes)
	switch {
	case err2 == io.EOF:
		fmt.Println("Reached EOF - close this connection.")
		return
	case err2 != nil:
		//fmt.Println("\nError reading command. Got: '"+line+"'\n", err)
		if err2 != nil {
			fmt.Println("Cannot write to connection.\n", err2)
		}
		return
	}
	// one piece of message received
	message := string(bodyBytes[:])
	handleStrings(message, msgWriter)

}
func handleStrings(cmd string, msgWriter *bufio.Writer) {
	fmt.Println("handle: " + cmd)
	longlongs := make([]byte, 999)
	for i:= range longlongs {
		longlongs[i] = 'b'
	}
	msgWriter.Write(longlongs)
	msgWriter.Flush()
}
