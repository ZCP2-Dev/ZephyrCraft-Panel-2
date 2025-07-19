package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
)

var isDebug bool

// 定义消息结构
type Message struct {
	Command string `json:"command"`
	Output  string `json:"output"`
	Error   string `json:"error"`
}

// WebSocket升级配置
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // 生产环境需限制来源
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERR]WebSocket升级失败: %v", err)
		return
	}
	defer conn.Close()
	var cmd *exec.Cmd
	// 启动控制台程序(默认为bedrock_server.exe)
	if isDebug {
		cmd = exec.Command(`D:\LiteZero_Project\ZephyrCraft-Panel-2\test\bedrock_server.exe`) //调试用路径
	} else {
		cmd = exec.Command(`.\bedrock_server.exe`)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		sendError(conn, "[ERR]启动进程失败: "+err.Error())
		return
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		sendError(conn, "[ERR]创建输入管道失败: "+err.Error())
		return
	}

	// 启动进程
	if err := cmd.Start(); err != nil {
		sendError(conn, "[ERR]启动进程失败: "+err.Error())
		return
	}
	defer cmd.Process.Kill() // 连接关闭时终止进程

	// 异步读取输出
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == nil { // 进程正常结束
					sendMessage(conn, Message{Output: "进程已结束"})
				} else {
					sendError(conn, "[ERR]读取输出失败: "+err.Error())
				}
				return
			}
			sendMessage(conn, Message{Output: line})
		}
	}()

	// 处理前端命令
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("[ERR]读取WebSocket消息失败: %v", err)
			return
		}

		// 发送命令到控制台
		cmdStr := msg.Command + "\n"
		_, err = stdin.Write([]byte(cmdStr))
		if err != nil {
			sendError(conn, "[ERR]发送命令失败: "+err.Error())
			return
		}
	}
}

// 发送消息到前端
func sendMessage(conn *websocket.Conn, msg Message) {
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("[ERR]发送消息失败: %v", err)
	}
}

func sendError(conn *websocket.Conn, errMsg string) {
	sendMessage(conn, Message{Error: errMsg})
}

func main() {
	isDebug = os.Getenv("DEBUG") == "1"
	if isDebug {
		log.Println("[ZephyCraft-Panel-2]现正于调试模式下运行")
	}
	// 注册WebSocket路由
	http.HandleFunc("/ws", handleWebSocket)

	// 启动HTTP服务
	log.Println("[ZephyCraft-Panel-2]websocket服务于端口1145上启动")
	log.Fatal(http.ListenAndServe(":1145", nil))
}
