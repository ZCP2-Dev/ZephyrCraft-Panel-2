package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var isDebug bool // 调试模式标记，用于控制是否输出调试信息

var config Config // 配置结构体实例，存储从配置文件读取的内容

type Config struct { // 定义配置结构，对应配置文件（如config.json）的格式
	Port       string `json:"port"`       // WebSocket 服务监听端口
	ServerPath string `json:"ServerPath"` // 要启动的服务程序路径
}

type Message struct { // 定义与前端交互的消息结构，前后端通过该结构传递命令、状态等数据
	Command string `json:"command"` // 命令类型，如 "start"、"stop"、"input"
	Content string `json:"content"` // 命令内容，比如输入的控制台指令
	Output  string `json:"output"`  // 服务程序输出内容
	Error   string `json:"error"`   // 错误信息，传递过程中出现错误时使用
	Status  string `json:"status"`  // 服务状态，如 "running"、"stopped"
}

// WebSocket 升级配置，用于将 HTTP 连接升级为 WebSocket 连接
var upgrader = websocket.Upgrader{
	// 生产环境建议严格校验来源，这里为了方便直接返回 true 允许所有来源
	CheckOrigin: func(r *http.Request) bool { return true },
}

// 读取配置文件函数，从指定路径加载配置并解析到 Config 结构体
func readConfig() Config {
	// 打开配置文件，路径为 .\Panel_Setting\config.json
	file, err := os.Open(".\\Panel_Setting\\config.json")
	if err != nil {
		// 打开失败时记录致命错误并终止程序
		log.Fatalf("[ERROR]无法打开配置文件: %v", err)
	}
	defer file.Close() // 函数结束前关闭文件，释放资源

	var config Config
	// 将文件内容解码到 Config 结构体实例
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		// 解码失败时记录致命错误并终止程序
		log.Fatalf("[ERROR]无法解析配置文件: %v", err)
	}
	return config
}

// 发送错误消息到前端函数，将错误信息封装成 Message 发送给前端
func sendError(conn *websocket.Conn, errMsg string) {
	// 尝试通过 WebSocket 连接发送错误消息
	if err := conn.WriteJSON(Message{Error: errMsg}); err != nil {
		// 发送失败时记录错误日志
		log.Printf("[ERROR]发送错误消息失败: %v", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// 升级失败时记录错误日志并返回
		log.Printf("[ERROR]WebSocket升级失败: %v", err)
		return
	}
	defer conn.Close() // 函数结束前关闭 WebSocket 连接

	// 创建进程管理器实例，传入 WebSocket 连接用于交互
	pm := getProcessManager()
	pm.UpdateConnection(conn)

	// 持续处理前端发送的消息
	for {
		var msg Message
		// 读取前端通过 WebSocket 发送的 JSON 消息并解析到 Message 结构体
		err := conn.ReadJSON(&msg)
		if err != nil {
			// 读取失败时记录错误日志并返回，终止循环
			if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
				log.Printf("[WARN]前端已主动退出连接") //如果是前端主动关闭连接，则不记录错误
			} else {
				log.Printf("[ERROR]读取WebSocket消息失败: %v", err) // 其他错误记录日志
			}
			pm.UpdateConnection(nil) // 清空连接
			return
		}

		// 根据消息的 Command 字段，分情况处理
		if isDebug {
			log.Printf("[Main][DEBUG]收到前端命令，类型: %s, 内容: %s", msg.Command, msg.Content) // 如果是调试模式，记录收到的命令和内容
		}
		switch msg.Command {
		case "start":
			// 处理启动命令，调用进程管理器的 StartProcess 方法
			if err := pm.StartProcess(config.ServerPath); err != nil {
				sendError(conn, "[ERROR]启动进程失败: "+err.Error())
			}
		case "stop":
			// 处理停止命令，调用进程管理器的 StopProcess 方法
			if err := pm.SendCommand("stop"); err != nil {
				sendError(conn, "[ERROR]停止进程失败: "+err.Error())
			}
		case "input":
			// 处理输入命令，调用进程管理器的 SendCommand 方法
			if err := pm.SendCommand(msg.Content); err != nil {
				sendError(conn, "[ERROR]发送命令失败: "+err.Error())
			}
		}
	}
}

var globalProcessManager *ProcessManager

func getProcessManager() *ProcessManager {
	if globalProcessManager == nil {
		globalProcessManager = NewProcessManager(nil)
	}
	return globalProcessManager
}
func main() {
	// 根据环境变量 DEBUG 是否为 "1"，设置调试模式标记
	isDebug = os.Getenv("DEBUG") == "1"
	if isDebug {
		log.Println("[ZephyCraft-Panel-2]现正于调试模式下运行")
	}
	config = readConfig() // 读取配置文件内容到 config 变量

	// 注册 WebSocket 路由，当访问 /ws 路径时，调用 handleWebSocket 函数处理
	http.HandleFunc("/ws", handleWebSocket)

	// 启动 HTTP 服务，监听配置文件中指定的端口
	log.Println("[ZephyCraft-Panel-2]websocket服务于端口" + config.Port + "上启动")
	// 启动服务，若失败则记录致命错误并终止程序（因 log.Fatal 会处理错误并退出）
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
