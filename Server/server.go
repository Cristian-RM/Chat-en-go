package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

const addr = "127.0.0.1:3000"
const bufferSize = 256
const endLine = 10

var clients []net.Conn

func main() {
	//Creando un slice de connecciones
	clients = make([]net.Conn, 0)
	//listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Can't listen on " + addr)
		os.Exit(1)
	}
	fmt.Println("Servidor en linea")
	for {
		conn, _ := listener.Accept()
		//agregar a un slice usamos append
		clients = append(clients, conn)
		fmt.Printf("Connecnion %s", strconv.Itoa(len(clients)))
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	//defer si se cierra el programa, o un return siempre se ejecuta.
	//finaly
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
				sendToOtherClients(conn, []byte("Se salio alguien?"))
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
		sendToOtherClients(conn, data)
		data = make([]byte, 0)
	}
}

func sendToOtherClients(sender net.Conn, data []byte) {
	if len(data) <= 0 {
		return
	}
	for i := 0; i < len(clients); i++ {
		if clients[i] != sender {
			clients[i].Write(data)
		}
	}

}
