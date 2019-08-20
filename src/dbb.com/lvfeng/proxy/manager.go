package proxy

import (
	"sync"
	"log"
)

var instance *ConnManger
var once sync.Once
var addconnchan = make(chan ConnectionPair, 1000)
var closeConnChan = make(chan ConnectionPair, 1000)
var cancelChan = make(chan struct{})
var monitorChan = make(chan InOutCount, 10000)
var lock sync.RWMutex
var watchD sync.WaitGroup

func GetInstance() *ConnManger{
	once.Do(func() {
		instance = & ConnManger{Connections:make(map[string]ConnectionPair), Started:false, Closing:false}
	})
	return instance
}

type InOutCount struct {
	in, out int
}

type ConnManger struct{
	Connections map[string]ConnectionPair
	Started bool
	Closing bool
	in int
	out int
}

func (m *ConnManger) AddConnection(connection ConnectionPair){
	addconnchan <- connection
}

func (m *ConnManger) ConnectionCount() int{
	return len(m.Connections)
}

func (m *ConnManger) Start(){
	if (*m).Started == true{
		return
	}
	(*m).Started = true
	go (*m).handleNewConnection()
	watchD.Add(1)
	go (*m).Monitor()
	watchD.Add(1)
	go (*m).ConnectionClose()
	watchD.Add(1)
}

func (m *ConnManger) Close(){
	defer watchD.Wait()
	if(*m).Closing == true{
		return
	}
	(*m).Closing = true
}


func (m *ConnManger)handleNewConnection(){
	for{
		select {
			case newConn := <- addconnchan:{
				lock.Lock()
				if _, ok := (*m).Connections[newConn.HostPort]; ok{
					// do nothing
				}
				(*m).Connections[newConn.HostPort] = newConn
				lock.Unlock()

			}
			case <- cancelChan:
				// Send close to connection
				watchD.Done()
				return
			//default:
			//	log.Print("HHAHAHA")

		}
	}
}

func ConnectionClose(conn ConnectionPair){
	closeConnChan <- conn
}

func TranCount(inc, outc int){
	monitorChan <- InOutCount{in: (int)(inc / 8), out:(int)(outc / 8)}
}

func (m *ConnManger)ConnectionClose(){
	for{
		select {
			case <- cancelChan:
				watchD.Done()
				return
			case conn:= <- closeConnChan:
				lock.Lock()
				var ok bool
				if _, ok = (*m).Connections[conn.HostPort]; ok{
					// do nothing
				}
				if ok{
					delete((*m).Connections, conn.HostPort)
					log.Printf("Close con: %s", conn.HostPort)
				}
				lock.Unlock()
		}
	}
}

func (m *ConnManger) Total() int{
	return (*m).in + (*m).out
}

// Monitor network traffic in byte.
func (m *ConnManger)Monitor(){
	for{
		select {
		case n := <- monitorChan:
			(*m).in += n.in
			(*m).out += n.out
			log.Printf("in+: %d, out+ %d, total: %d, in:%d, out:%d", n.in, n.out, (*m).Total(), (*m).in,
				(*m).out)
		case <- cancelChan:
			// Send close to connection
			watchD.Done()
			return
		}
	}
}