package proxy
import (
	"io"
	"log"
	"net"
	"runtime/debug"
	"time"
)

type Command int

const (
	CommandClose Command = iota
	CommandReconnect
)

type Addr struct {
	host string
	port int
}

type NIOMonitorInterface interface{
	Read(b []byte) (n int, err error)

	Write(b []byte) (n int, err error)

	Close() error
}


const TargetHost string = "127.0.0.1"
const TargetPort int = 8080

type ConnectionPair struct{
	InByteCount int
	OutByteCount int
	HostPort string
	InConn *net.Conn
	OutConn *net.Conn
	LastUseAt int64
	Enabled bool
	Closing bool
	CommandChan chan int
	DoneChan chan bool
}


func (conn *ConnectionPair)Enable() {
	if (*conn).Closing{
		return
	}
	(*conn).Enabled = true
	// start new goroutine
	//  1. record IO traffic, set alive
	//  2. monitor connection, e.g. close(reconnect),
}

func (c *ConnectionPair) SetOutConnection(outCon *net.Conn){
	(*c).OutConn = outCon
}

func monitorIO(conn *ConnectionPair) {
	for {
		select {
		case command := <-(*conn).CommandChan:
			{
				if command == int(CommandClose) {
					(*conn).Closing = true
				}

				if command == int(CommandReconnect) {
					//TODO
				}
			}
		default:
		}

	}
}

func (c ConnectionPair) IsEnabled() bool{
	if (c).Closing{
		return false
	}
	return c.Enabled
}

func (conn *ConnectionPair) Disable(){
	(*conn).Enabled = false
}

func (conn *ConnectionPair) Update(){
	(*conn).LastUseAt = time.Now().UTC().UnixNano()
}

func (conn ConnectionPair) TimeOut() bool{
	if time.Now().UTC().UnixNano() - conn_max_idle_time < conn.LastUseAt{
		return false
	}
	return true
}

func (c ConnectionPair) Trans(count int, dir Direction){
	switch dir {
	case TransDirectionUp:
		c.InByteCount += count
		TranCount(count,0)
		break
	case TransDirectionDown:
		c.OutByteCount += count
		TranCount(0, count)
		break
	}
}

type Direction int

type TransResult struct{
	direction Direction
	written int
	err error
	suc bool
}

const (
	TransDirectionUp = iota
	TransDirectionDown
)

func (c *ConnectionPair)TransferIO(){
	result := make(chan TransResult, 2)
	go transferData(*c, TransDirectionUp, result)
	go transferData(*c, TransDirectionDown, result)
	go cleaner(*c, result)
	return
}

func transferData(conn ConnectionPair, dir Direction, result chan TransResult){
	var err error
	var count int
	var suc bool = true
	defer func()  {
		if e := recover(); e != nil {
			log.Printf("transfer IO crased, err: %s\n, trace: %s", e, string(debug.Stack()))
		}
		result <- TransResult{written:count, err: err, suc: suc, direction:dir}
	}()

	//var suc bool
	buf := make([]byte, 2*1024)
	var src, dst net.Conn
	switch dir {
	case TransDirectionUp:
		src, dst = *conn.InConn, *conn.OutConn
		break
	case TransDirectionDown:
		src, dst = *conn.OutConn, *conn.InConn
		break
	default:
		return
	}

	for{
		select {
		case <- conn.DoneChan:
			return
		default:
			nr, er := copyData(src, dst, buf)
			count += nr
			if er != nil{
				err = er
				suc = false
				return
			}
			conn.Trans(nr, dir)
		}
	}
}


func copyData(src, dst io.ReadWriter, buf []byte)(int, error){
	var count int
	var err error
	//buf := make([]byte, 1*1024)

	nr, er := src.Read(buf)
	if nr >0{
		nw, ew := dst.Write(buf[0:nr])
		count += nw
		if ew != nil{
			err = ew
		}
		if nw != nr{
			err = ew
		}
	}
	if er != nil{
		err = er
	}
	return nr, err
}



func cleaner(conn ConnectionPair, result chan TransResult){
	//defer func(){
	//	src.Close()
	//	dst.Close()
	//}()
	src := io.ReadWriteCloser(*conn.InConn)
	dst := io.ReadWriteCloser(*conn.OutConn)
	var count int
	for{
		select {
		case r := <- result:
			if r.direction == TransDirectionUp{
				src.Close()
				dst.Close()
				count += 1
				log.Printf("Upstream closed.")
			}
			if r.direction == TransDirectionDown{
				src.Close()
				dst.Close()
				log.Printf("Downstream closed.")
				count += 1
			}
			if count == 2 {
				log.Print("Up&Donw sream pair closed..")
				ConnectionClose(conn)
				return
			}
		}
	}
}