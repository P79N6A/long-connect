package comet

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)
type Data struct {
	OwnId string
	SendTo string
	Msg string
}

var (
	ConnInfoList ConnList
	lock sync.RWMutex
)
func init() {
	ConnInfoList.Conn = make(map[string]*NetConn)
}
func StartServer(network, port string) {
	listen, err := net.Listen(network, port)
	if err != nil {
		panic(fmt.Errorf("StartServer error: %v", err))
	}
	fmt.Printf("Server Start...\n")
	for {
		data := make([]byte, 1024)
		conn ,err := listen.Accept()
		if err != nil {
			printErrorLog(err)
			continue
		}
		size, err := conn.Read(data)
		if err != nil {
			printErrorLog(err)
			continue
		}
		var jsonData Data
		err = json.Unmarshal(data[:size], &jsonData)
		//fmt.Println(string(data[:size]))
		if err != nil {
			printErrorLog(err)
			continue
		}
		netConn := &NetConn{
			ConnInfo: conn,
			Status: true,
		}
		lock.Lock()
		err = ConnInfoList.Set(jsonData.OwnId, netConn)
		lock.Unlock()
		// 发送消息
		go Message(conn, jsonData.OwnId)
		//
	}
}
func Message(connFrom net.Conn, key string) {
	for {
		fromData := make([]byte, 1024)
		size, err := connFrom.Read(fromData)
		if err != nil {
			if err == io.EOF {
				lock.Lock()
				delete(ConnInfoList.Conn, key)
				lock.Unlock()
				break
			}
			printErrorLog(err)
			continue
		}
		var jsonData Data
		err = json.Unmarshal(fromData[:size], &jsonData)
		//fmt.Println(string(fromData[:size]))
		if err != nil {
			printErrorLog(err)
			continue
		}
		//key := jsonData.SendTo
		lock.RLock()
		conn := ConnInfoList.Conn[key].ConnInfo
		lock.RUnlock()
		conn.Write(fromData[:size])

		//select {
		//case <-ConnInfoList.Conn[key].StopChan :
		//	break
		//}
	}
}
func ClearConnList() {
	for key, list := range ConnInfoList.Conn {
		if !list.Status {
			delete(ConnInfoList.Conn, key)
			// 发送信号，停止线程
			//ConnInfoList.Conn[key].StopChan <- true
		}
	}
	// 间隔一段时间清理一次长连接
	time.Sleep(time.Minute * 1)
}

func CurrentLongConnection() {
	for {
		fmt.Printf("Current Connection num is: %d\n", len(ConnInfoList.Conn))
		time.Sleep(time.Second * 5)
	}
}
func printErrorLog(err error) {
	fmt.Fprintln(os.Stderr, err)
}