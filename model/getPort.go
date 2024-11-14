package model
import (
	"net"
)
// GetFreePort 获取一个可用的 TCP 端口
func GetFreePort() (int, error) {
	// 解析 TCP 地址，使用 "localhost:0" 表示系统自动分配一个可用端口
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	// 监听 TCP 地址，获取一个可用端口
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	// 返回监听地址的端口号
	return l.Addr().(*net.TCPAddr).Port, nil
}
