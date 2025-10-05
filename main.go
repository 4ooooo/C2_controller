package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	// serverAddr = "101.200.236.51:8080" // C2服务器地址，根据需要修改
	serverAddr = "127.0.0.1:8080" // 本地测试地址
)

func main() {
	fmt.Println("C2控制端启动中...")
	fmt.Printf("正在连接到C2服务器 %s...\n", serverAddr)

	// 连接到C2服务器
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("连接服务器失败: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 发送控制端标识符
	_, err = conn.Write([]byte{'C'})
	if err != nil {
		fmt.Printf("发送控制端标识符失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("已连接到C2服务器")

	// 创建用于从服务器接收数据的goroutine
	go receiveFromServer(conn)

	// 处理用户输入并发送到服务器
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()

		if cmd == "exit" {
			fmt.Println("正在断开连接...")
			break
		}

		// 发送命令到服务器
		_, err = conn.Write([]byte(cmd + "\n"))
		if err != nil {
			fmt.Printf("发送命令失败: %v\n", err)
			break
		}
	}

	if err = scanner.Err(); err != nil {
		fmt.Printf("读取输入失败: %v\n", err)
	}
}

// receiveFromServer 接收并显示服务器消息
func receiveFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("从服务器接收消息失败: %v\n", err)
			os.Exit(1)
		}

		// 处理命令提示符
		if msg == "> \n" {
			fmt.Print("> ")
			continue
		}

		// 直接打印服务器消息，不进行编码转换
		// 因为服务器通过DNS协议传输的数据已经是正确的UTF-8编码
		fmt.Print(msg)
	}
}
