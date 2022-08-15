package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

const addr = "127.0.0.1:3000"
const bufferSize = 256
const endLine = 10

var nick string

//es un reader de la consola
// video numero 6 https://www.youtube.com/watch?v=GPXz_9BuE8Y
var in *bufio.Reader

func main() {
	//entrada estandar del sistema operativo
	in = bufio.NewReader(os.Stdin)
	for nick == "" {
		fmt.Printf("Dame tu nick: \n")
		buf, _, _ := in.ReadLine()
		nick = string(buf)
	}
	var conn net.Conn
	var err error
	for {
		fmt.Printf("Conectando a %s... \n", addr)
		conn, err = net.Dial("tcp", addr)
		if err == nil {
			break
		}
	}
	defer conn.Close()

	go reciveMessages(conn)
	handleConnection(conn)

}

func handleConnection(conn net.Conn) {

	conn.Write(append([]byte(nick + " se ha conectado. \n")))
	for {
		for {
			buf, _, _ := in.ReadLine()
			if len(buf) > 0 {
				conn.Write(append([]byte(nick+" ->"), append(buf, endLine)...))
			}
		}
	}
}

func reciveMessages(conn net.Conn) {
	defer conn.Close()
	var data []byte
	buffer := make([]byte, bufferSize)

	for {
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
			}
			if n <= 0 {
				break
			}
			//replace, los campos vacíos en hexadecimal son \x00
			//y Trim retorna un subslice que no contiene el parametro pasado
			//por lo que este paso elimina los espacios vácios
			buffer = bytes.Trim(buffer[:n], "\x00")
			data = append(data, buffer...)
			if data[len(data)-1] == endLine {
				break
			}
		}
		fmt.Printf("%s\n", data[:len(data)-1])
		data = make([]byte, 0)
	}
}
