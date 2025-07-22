package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/gorilla/websocket"
)

// ProcessManager 管理外部进程的启动和通信
type ProcessManager struct {
	cmd     *exec.Cmd
	stdout  io.ReadCloser
	stdin   io.WriteCloser
	running bool
	mu      sync.Mutex      // 用于并发控制
	conn    *websocket.Conn // WebSocket连接
}

// NewProcessManager 创建新的进程管理器
func NewProcessManager(conn *websocket.Conn) *ProcessManager {
	return &ProcessManager{
		conn:    conn,
		running: false,
	}
}

// StartProcess 启动外部进程
func (pm *ProcessManager) StartProcess(serverPath string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.running {
		return nil // 进程已在运行
	}

	// 创建命令
	pm.cmd = exec.Command(".\\Panel_Setting\\pty-proxy.exe", serverPath)

	var err error
	// 设置输入输出管道
	if pm.stdout, err = pm.cmd.StdoutPipe(); err != nil {
		return err
	}

	if pm.stdin, err = pm.cmd.StdinPipe(); err != nil {
		return err
	}

	// 启动进程
	if err := pm.cmd.Start(); err != nil {
		return err
	}

	pm.running = true

	// 异步读取输出
	go pm.readOutput()

	return nil
}

// StopProcess 停止外部进程
func (pm *ProcessManager) StopProcess() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if !pm.running || pm.cmd == nil {
		return nil // 进程未运行
	}

	if err := pm.cmd.Process.Kill(); err != nil {
		return err
	}

	pm.running = false
	return nil
}

// SendCommand 向进程发送命令
func (pm *ProcessManager) SendCommand(command string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if !pm.running || pm.stdin == nil {
		return nil // 进程未运行
	}

	// 发送命令到控制台
	cmdStr := command + "\r\n"
	_, err := pm.stdin.Write([]byte(cmdStr))
	return err
}

// readOutput 异步读取进程输出并发送到WebSocket
func (pm *ProcessManager) readOutput() {
	reader := bufio.NewReader(pm.stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			pm.StopProcess() // 发生错误时停止进程
			return
		}

		// 发送输出到前端
		pm.sendMessage(Message{Output: line})
	}
}

// sendMessage 发送消息到前端
func (pm *ProcessManager) sendMessage(msg Message) {
	if pm.conn == nil {
		return
	}

	if err := pm.conn.WriteJSON(msg); err != nil {
		log.Printf("[ERR]发送消息失败: %v", err)
	}
}
