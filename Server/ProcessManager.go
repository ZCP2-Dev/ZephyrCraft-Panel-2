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
// 负责控制外部程序（比如通过 pty-proxy.exe 启动的服务）的启停、命令交互，
// 以及收集程序输出并通过 WebSocket 发给前端
type ProcessManager struct {
	cmd     *exec.Cmd       // 用于管理外部进程的命令对象，包含进程启动参数、运行状态等
	stdout  io.ReadCloser   // 外部进程的标准输出管道，用于读取进程输出内容
	stdin   io.WriteCloser  // 外部进程的标准输入管道，用于向前端进程发送命令
	running bool            // 标记进程是否正在运行
	mu      sync.Mutex      // 互斥锁，用于并发场景下保护进程相关状态（比如 running、cmd 等），避免竞争条件
	conn    *websocket.Conn // 与前端通信的 WebSocket 连接，用于发送输出、状态等消息
}

// NewProcessManager 创建新的进程管理器
// 初始化时传入 WebSocket 连接，后续用它给前端发消息
func NewProcessManager(conn *websocket.Conn) *ProcessManager {
	return &ProcessManager{
		conn:    conn,  // 关联传入的 WebSocket 连接
		running: false, // 初始时进程未运行
	}
}

// StartProcess 启动外部进程
// serverPath 是要启动的目标程序路径（通过 pty-proxy.exe 代理启动）
func (pm *ProcessManager) StartProcess(serverPath string) error {
	pm.mu.Lock()         // 加锁，避免并发操作进程状态
	defer pm.mu.Unlock() // 函数结束后自动解锁

	if pm.running {
		return nil // 进程已经在运行，直接返回，避免重复启动
	}

	// 创建命令对象，指定要执行的程序（pty-proxy.exe）和参数（serverPath 即实际要启动的服务路径）
	pm.cmd = exec.Command(".\\Panel_Setting\\pty-proxy.exe", serverPath)

	var err error
	// 建立标准输出管道，后续通过它读取进程输出
	if pm.stdout, err = pm.cmd.StdoutPipe(); err != nil {
		return err // 管道建立失败，返回错误
	}

	// 建立标准输入管道，后续通过它给进程发命令
	if pm.stdin, err = pm.cmd.StdinPipe(); err != nil {
		return err // 管道建立失败，返回错误
	}

	// 真正启动外部进程
	if err := pm.cmd.Start(); err != nil {
		return err // 进程启动失败，返回错误
	}

	pm.running = true // 标记进程已启动运行

	// 异步启动一个 goroutine 读取进程输出，避免阻塞当前逻辑
	go pm.readOutput()

	return nil // 启动成功，返回 nil
}

// StopProcess 停止外部进程
func (pm *ProcessManager) StopProcess() error {
	pm.mu.Lock()         // 加锁，保证操作进程状态时的并发安全
	defer pm.mu.Unlock() // 函数结束自动解锁

	// 进程未运行或者命令对象为空，直接返回
	if !pm.running || pm.cmd == nil {
		return nil
	}

	// 尝试杀死进程（强制终止）
	if err := pm.cmd.Process.Kill(); err != nil {
		return err // 终止失败，返回错误
	}

	pm.running = false // 标记进程已停止
	return nil         // 停止成功，返回 nil
}

// SendCommand 向进程发送命令
// command 是前端传来的要执行的指令（比如服务控制台命令）
func (pm *ProcessManager) SendCommand(command string) error {
	pm.mu.Lock()         // 加锁，避免并发写标准输入
	defer pm.mu.Unlock() // 函数结束自动解锁
	if isDebug {
		log.Printf("[ProcessManager][DEBUG]发送命令到控制台: %s", command) // 如果是调试模式，记录发送的命令
	}
	// 进程未运行或者标准输入管道为空，无法发命令，直接返回
	if !pm.running || pm.stdin == nil {
		return nil
	}

	// 拼接命令（加换行符模拟控制台输入回车），转成字节切片写入标准输入管道
	cmdStr := command + "\r\n"
	_, err := pm.stdin.Write([]byte(cmdStr))
	return err // 返回写入是否成功的错误
}

// readOutput 异步读取进程输出并发送到 WebSocket
// 持续从标准输出读内容，读到后通过 WebSocket 发给前端
func (pm *ProcessManager) readOutput() {
	// 用带缓冲的 reader 读取标准输出，按行处理
	reader := bufio.NewReader(pm.stdout)
	for {
		// 尝试读取一行内容（以换行符结尾）
		line, err := reader.ReadString('\n')
		if err != nil {
			// 读取出错（比如进程退出、管道断开等），停止进程并返回
			pm.StopProcess()
			return
		}

		// 将读取到的一行输出，封装成 Message 结构体，通过 WebSocket 发给前端
		pm.sendMessage(Message{Output: line})
	}
}

// sendMessage 发送消息到前端
// msg 是要发送的消息内容（包含输出、状态等信息）
func (pm *ProcessManager) sendMessage(msg Message) {
	// WebSocket 连接为空，无法发送，直接返回
	if pm.conn == nil {
		return
	}
	if isDebug {
		log.Printf("[ProcessManager][DEBUG]发送消息到前端: %s", msg) // 如果是调试模式，记录发送的消息内容
	}
	// 尝试通过 WebSocket 连接发送 JSON 格式消息
	if err := pm.conn.WriteJSON(msg); err != nil {
		// 发送失败，记录错误日志
		log.Printf("[ERR]发送消息失败: %v", err)
	}
}
