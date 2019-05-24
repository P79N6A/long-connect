package comet

import (
	"fmt"
	"net"
)

type NetConn struct {
	ConnInfo net.Conn // 连接Conn
	Status bool		// 连接状态
	StopChan chan bool // 线程状态
}
type ConnList struct {
	Conn map[string]*NetConn
}
func (n *ConnList) SetConnStatus(key string, status bool) error {
	if _, ok := n.Conn[key]; !ok {
		return fmt.Errorf("%s not in ConnList: %v\n", key, ok)
	}
	n.Conn[key].Status = status
	return nil
}

func (n *ConnList) Set(key string, netConn *NetConn) error {
	n.Conn[key] = netConn
	return nil
}

func (n *ConnList) Del(key string) {
	delete(n.Conn, key)
}

func (n *ConnList) GetAll() []string {
	data := make([]string, 0)
	for key, _ := range n.Conn {
		data = append(data, key)
	}
	return data
}
