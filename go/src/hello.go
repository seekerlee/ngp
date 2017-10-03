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
	fmt.Printf("hello, world\n")
	// Listen on a random endpoint.
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

func write2js(msgWriter *bufio.Writer, go2jsMessage <-chan string) {
	for {
		msg := <-go2jsMessage
		fmt.Println("write to js: " + msg)
		msgWriter.Write([]byte(msg))
		msgWriter.Flush()
	}
}

func byteArr2Uint(p []byte) uint32 {
    data := binary.BigEndian.Uint32(p)
	fmt.Println(data)
	return data
}
func uint2ByteArr(val uint32) []byte {
	bs := make([]byte, 4)
    binary.BigEndian.PutUint32(bs, val)
	fmt.Println(bs)
	return bs
}
// func main() {
// 	fmt.Println("Hello, playground")
// 	sss := sizeByteToSize([]byte{byte(0),byte(0),byte(3),byte(231)})
// 	fmt.Println(sss)
// }
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

readLoop:
	for {
		//line, err := msgReader.ReadString('\n')
		sizebytes := make([]byte, 4)
		line, err := msgReader.Read(sizebytes)
		switch {
		case err == io.EOF:
			fmt.Println("Reached EOF - close this connection.")
			break readLoop
		case err != nil:
			//fmt.Println("\nError reading command. Got: '"+line+"'\n", err)
			if err != nil {
				fmt.Println("Cannot write to connection.\n", err)
			}
			break readLoop
		}
		fmt.Println(line)
		bodyLen := byteArr2Uint(sizebytes)
		bodyBytes := make([]byte, bodyLen)
		_, err2 := msgReader.Read(bodyBytes)
		switch {
		case err2 == io.EOF:
			fmt.Println("Reached EOF - close this connection.")
			break readLoop
		case err2 != nil:
			//fmt.Println("\nError reading command. Got: '"+line+"'\n", err)
			if err2 != nil {
				fmt.Println("Cannot write to connection.\n", err2)
			}
			break readLoop
		}
		// one piece of message received
		message := string(bodyBytes[:])
		handleStrings(message, go2jsMessage)
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
