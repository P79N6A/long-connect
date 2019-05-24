package main

import (
	"bufio"
	"client/message"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)
var (
	lock sync.Mutex
	remoteAddr = &net.TCPAddr{}
)
func init() {
	var err error
	remoteAddr, err = net.ResolveTCPAddr("tcp", "yq01-abc-im-dev.epc.baidu.com:9999")
	if err != nil {
		panic(err)
	}
}
func main() {
	// 测试虚拟id地址
	localIPs := []string{
		"10.94.35.100",
		"10.94.35.101",
		"10.94.35.102",
		"10.94.35.103",
		"10.94.35.104",
		"10.94.35.105",
		"10.94.35.106",
		"10.94.35.107",
		"10.94.35.109",
		"10.94.35.110",
		"10.94.35.111",
		"10.94.35.112",
		"10.94.35.113",
		"10.94.35.114",
		"10.94.35.115",
		"10.94.35.116",
		"10.94.35.117",
		"10.94.35.118",
		"10.94.35.119",
		"10.94.35.120",
	}
	fmt.Println("please input max client num and sleep time(millisecond)")
	count := int64(0)
	// 虚拟网卡端口从10000开始
	startPort := 10000
	for {
		fmt.Printf("example: 10000 1: ")
		var maxThread, timeSleep int64
		// 输入本次要启动的线程数和时间间隔(实际线程数=输入的线程数*虚拟网卡个数)
		// 输入2000 20，虚拟网卡数：20，实际线程数=2000 * 20
		fmt.Scanf("%d%d", &maxThread, &timeSleep)
		for i := int64(0); i <= maxThread; i++ {
			for _, virtualIP := range localIPs {
				portStr := strconv.Itoa(startPort)
				localAddr, _:= net.ResolveTCPAddr("tcp", virtualIP + ":" + portStr)
				clientId := string(localAddr.IP) + ":" + strconv.Itoa(localAddr.Port)
				go StartClient(clientId, localAddr)
				count++
			}
			startPort ++
			time.Sleep(time.Millisecond * time.Duration(timeSleep))
			// 防止并发太高
			if i % 2000 == 0 {
				fmt.Printf("累计连接数：%d\n", count)
				time.Sleep(time.Second * 2)
			}
		}
		fmt.Printf("累计连接数：%d\n", count)
	}
	time.Sleep(time.Second * 100000)

}
func StartClient(threadId string, localAddr *net.TCPAddr) {
	conn, err := net.DialTCP("tcp", localAddr, remoteAddr)
	if err != nil {
		printErrorLog(err)
		return
	}
	defer conn.Close()
	var sendData message.Data
	sendData.SendTo = threadId
	sendData.OwnId = threadId
	sendData.Msg = "Hello Go"
	byteData, err := json.Marshal(sendData)
	//fmt.Println(string(byteData))
	if err != nil {
		printErrorLog(err)
		return
	}
	_, err = conn.Write(byteData)
	if err != nil {
		printErrorLog(err)
		return
	}
	// 启动获取消息线程
	go ReadData(conn)
	for {
		time.Sleep(time.Second * 20)
		//line , _, err := reader.ReadLine()
		//if err != nil {
		//	fmt.Printf("Read data error: %v\n", err)
		//}
		//input, _ := strconv.Atoi(string(line))
		//if input == 1 {
		//	fmt.Printf("your id is: %d\n", sendData.OwnId)
		//	continue
		//} else if input == 4 {
		//	break
		//} else if input == 3 {
		//	var send int64
		//	fmt.Printf("input friend id: ")
		//	fmt.Scanf("%d", &send)
		//	sendData.SendTo = send
		//	continue
		//}
		//sendData.Msg = string(line)
		var sendData message.Data
		sendData.SendTo = threadId
		sendData.OwnId = threadId
		sendData.Msg = "Hello Go"
		tmp, err := json.Marshal(sendData)
		if err != nil {
			printErrorLog(err)
			continue
		}
		//fmt.Println(string(tmp))
		writeBuf := bufio.NewWriter(conn)
		_, err = writeBuf.Write(tmp)
		writeBuf.Flush()
		if err != nil {
			printErrorLog(err)
		}
		//fmt.Printf("\t\t you said: %s\n", line)
	}
}
func ReadData(conn net.Conn) {
	defer conn.Close()
	jsonData := make([]byte, 1024)
	for {
		_, err := conn.Read(jsonData)
		if err != nil {
			if err == io.EOF {
				break
			}
			printErrorLog(err)
			continue
		}
		//fmt.Println(jsonData[:size])
		//var data message.Data
		//err = json.Unmarshal(jsonData[:size], &data)
		//if err != nil {
		//	fmt.Printf("Conn Read Data unmarshal error: %v\n", err)
		//	continue
		//}
		//fmt.Printf("\n\t\t%d Say: %s\n", data.OwnId, data.Msg)
		//fmt.Printf("%s", "Please input message: ")
	}
}
func printErrorLog(err error) {
	fmt.Fprintln(os.Stderr, err)
}