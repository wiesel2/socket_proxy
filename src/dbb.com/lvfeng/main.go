package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"socket_proxy/src/dbb.com/lvfeng/cfg"
	"socket_proxy/src/dbb.com/lvfeng/proxy"
	"strconv"
)
//
//
//
func trans(conn *net.Conn, m *proxy.ConnManger){
	log.Printf("New Connect connected: %v", (*conn).RemoteAddr())

	outConnection, err :=  net.Dial("tcp", proxy.TargetHost + ":" + strconv.Itoa(proxy.TargetPort))
	if err != nil{
		log.Fatal(err)
	}
	newConn := proxy.ConnectionPair{
		InConn:conn,
		OutConn:&outConnection,
		CommandChan:make(chan int, 100),
		HostPort:(*conn).RemoteAddr().String(),
		DoneChan: make(chan bool),
	}
	(*m).AddConnection(newConn)
	newConn.TransferIO()
}

func main() {
	log.Printf("Socket proxy launching, pid: %d", os.Getpid())
	localAddr := fmt.Sprintf("%s:%d", cfg.DefaultCfg.ServerConfig.LocalHost, cfg.DefaultCfg.ServerConfig.LocalPort)
	listener, err := net.Listen("tcp", localAddr)
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
