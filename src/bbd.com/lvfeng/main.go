package main

import (
	"net"
	"log"
	"bbd.com/lvfeng/proxy"
	"strconv"
	"os"
)

func trans(conn *net.Conn, m *proxy.ConnManger){
	log.Printf("New Connect connected: %v", (*conn).RemoteAddr())

	outConnection, err :=  net.Dial("tcp", proxy.TargetHost + ":" + strconv.Itoa(proxy.TargetPort))
	if err != nil{
		log.Fatal(err)
	}
	new_conn := proxy.ConnectionPair{
		InConn:conn,
		OutConn:&outConnection,
		CommandChan:make(chan int, 100),
		HostPort:(*conn).RemoteAddr().String(),
		DoneChan: make(chan bool),
	}
	(*m).AddConnection(new_conn)
	new_conn.TransferIO()

}


func main() {
	log.Printf("Socket proxy launching, %d", os.Getpid())
	listener, err := net.Listen("tcp", "localhost:8001")
	m := proxy.GetInstance()
	m.Start()
	if err != nil{
		log.Fatal(err)
	}
	// TODO| create a another goroutine to accept new connection
	// TODO|

	for {
		conn, err := listener.Accept()
		if err != nil{
			log.Print(err)
			continue
		}

		go trans(&conn, m)
		}
	}
