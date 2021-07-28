package snake

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	slog "github.com/m2c/kiplestar/commons/log"
	_ "math/rand"
	"net"
	"sync"
)

var node *snowflake.Node

var rwlock sync.RWMutex

func InitSnake() {
	ip, iperror := Lower16BitPrivateIP()
	if iperror != nil {
		slog.Error(iperror.Error())
		panic("No IP address assigned...")
	}
	fmt.Println(ip)
	newNode, nodeError := snowflake.NewNode(int64(ip))
	if iperror != nil {
		slog.Info(nodeError.Error())
		panic("create the node is failed by the ip...")
	}
	node = newNode
}

func GetSnokeNode() *snowflake.Node {
	if node == nil {
		rwlock.Lock()
		defer rwlock.Unlock()
		if node != nil {
			return node
		}
		InitSnake()
	}
	return node
}

func Lower16BitPrivateIP() (uint16, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}
