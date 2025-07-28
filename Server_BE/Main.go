package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"encoding/base64"
	"strings"

	"github.com/gorilla/websocket"
)

var isDebug bool // 调试模式标记，用于控制是否输出调试信息

var config Config // 配置结构体实例，存储从配置文件读取的内容

type Config struct { // 定义配置结构，对应配置文件（如config.json）的格式
	Port          string `json:"port"`                    // WebSocket 服务监听端口
	ServerPath    string `json:"ServerPath"`              // 要启动的服务程序路径
	Version       string `json:"Version,omitempty"`       // 服务器版本号
	LoaderVersion string `json:"LoaderVersion,omitempty"` // 加载器版本号
	Uniteban      bool   `json:"Uniteban"`                // 云黑校验开关，默认true
}

type Message struct { // 定义与前端交互的消息结构，前后端通过该结构传递命令、状态等数据
	Command string `json:"command"` // 命令类型，如 "start"、"stop"、"input"
	Content string `json:"content"` // 命令内容，比如输入的控制台指令
	Output  string `json:"output"`  // 服务程序输出内容
	Error   string `json:"error"`   // 错误信息，传递过程中出现错误时使用
	Status  string `json:"status"`  // 服务状态，如 "running"、"stopped"

	// 新增字段：系统状态信息
	SystemInfo *SystemInfo `json:"systemInfo,omitempty"` // 系统状态信息
	ServerInfo *ServerInfo `json:"serverInfo,omitempty"` // 服务器信息
	Players    []*Player   `json:"players,omitempty"`    // 玩家列表

	// 文件管理相关字段
	FilePath    string     `json:"filePath,omitempty"`    // 文件路径
	FileContent string     `json:"fileContent,omitempty"` // 文件内容
	FileList    []FileInfo `json:"fileList,omitempty"`    // 文件列表
	OldPath     string     `json:"oldPath,omitempty"`     // 旧路径（重命名用）
	NewPath     string     `json:"newPath,omitempty"`     // 新路径（重命名用）

	// 压缩相关字段
	FilesToZip  []string `json:"filesToZip,omitempty"`  // 要压缩的文件列表
	ZipFileName string   `json:"zipFileName,omitempty"` // 压缩文件名
}

// SystemInfo 系统状态信息
type SystemInfo struct {
	CPUUsage    float64 `json:"cpuUsage"`    // CPU使用率
	MemoryUsage float64 `json:"memoryUsage"` // 内存使用率
	MemoryTotal uint64  `json:"memoryTotal"` // 总内存(MB)
	MemoryUsed  uint64  `json:"memoryUsed"`  // 已用内存(MB)
	DiskUsage   float64 `json:"diskUsage"`   // 磁盘使用率
	DiskTotal   uint64  `json:"diskTotal"`   // 总磁盘空间(MB)
	DiskUsed    uint64  `json:"diskUsed"`    // 已用磁盘空间(MB)
	Uptime      uint64  `json:"uptime"`      // 系统运行时间(秒)
}

// ServerInfo 服务器信息
type ServerInfo struct {
	Version       string `json:"version"`                 // 服务器版本
	LoaderVersion string `json:"loaderVersion,omitempty"` // 加载器版本
	StartTime     string `json:"startTime"`               // 启动时间
	PlayerCount   int    `json:"playerCount"`             // 在线玩家数
	MaxPlayers    int    `json:"maxPlayers"`              // 最大玩家数
	Uptime        uint64 `json:"uptime"`                  // 服务器运行时间(秒)
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
	// 默认值处理
	if config.Version == "" {
		config.Version = ""
	}
	if config.LoaderVersion == "" {
		config.LoaderVersion = ""
	}
	// Uniteban默认开启
	// 注意：json反序列化bool为false时，只有字段缺失才会为false
	if _, ok := getRawField(".\\Panel_Setting\\config.json", "Uniteban"); !ok {
		config.Uniteban = true
	}
	return config
}

// 获取json原始字段是否存在
func getRawField(path, key string) (interface{}, bool) {
	file, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer file.Close()
	var m map[string]interface{}
	if err := json.NewDecoder(file).Decode(&m); err != nil {
		return nil, false
	}
	v, ok := m[key]
	return v, ok
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
			if err := pm.StopProcess(); err != nil {
				sendError(conn, "[ERROR]停止进程失败: "+err.Error())
			}
		case "input":
			// 处理输入命令，调用进程管理器的 SendCommand 方法
			if err := pm.SendCommand(msg.Content); err != nil {
				sendError(conn, "[ERROR]发送命令失败: "+err.Error())
			}
		case "status":
			// 处理状态查询命令，返回当前进程状态
			status := "stopped"
			if pm.IsRunning() {
				status = "running"
			}
			// 直接发送到WebSocket，不通过ProcessManager
			if err := conn.WriteJSON(Message{Status: status}); err != nil {
				log.Printf("[ERROR]发送状态消息失败: %v", err)
			}
		case "getSystemInfo":
			// 处理系统信息查询命令
			systemMonitor := GetSystemMonitor()
			if systemMonitor != nil {
				systemInfo := systemMonitor.GetSystemInfo()
				// 直接发送到WebSocket，不通过ProcessManager
				if err := conn.WriteJSON(Message{SystemInfo: systemInfo}); err != nil {
					log.Printf("[ERROR]发送系统信息失败: %v", err)
				}
			}
		case "getServerInfo":
			// 处理服务器信息查询命令
			systemMonitor := GetSystemMonitor()
			if systemMonitor != nil {
				serverInfo := systemMonitor.GetServerInfo()
				// 直接发送到WebSocket，不通过ProcessManager
				if err := conn.WriteJSON(Message{ServerInfo: serverInfo}); err != nil {
					log.Printf("[ERROR]发送服务器信息失败: %v", err)
				}
			}
		case "getPlayers":
			// 处理玩家列表查询命令
			playerManager := GetPlayerManager()
			players := playerManager.GetPlayers()
			// 直接发送到WebSocket连接
			if err := conn.WriteJSON(Message{Players: players}); err != nil {
				log.Printf("[ERROR]发送玩家列表失败: %v", err)
			}
			if isDebug {
				log.Printf("[Main][DEBUG]发送玩家列表，玩家数量: %d", len(players))
			}

		// 文件管理相关命令
		case "getFileList":
			// 获取文件列表
			fileManager := GetFileManager()
			files, err := fileManager.GetFileList(msg.Content)
			if err != nil {
				sendError(conn, "[ERROR]获取文件列表失败: "+err.Error())
			} else {
				if err := conn.WriteJSON(Message{FileList: files}); err != nil {
					log.Printf("[ERROR]发送文件列表失败: %v", err)
				}
			}

		case "readFile":
			// 读取文件内容
			fileManager := GetFileManager()

			content, err := fileManager.ReadFile(msg.FilePath)
			if err != nil {
				sendError(conn, "[ERROR]读取文件失败: "+err.Error())
			} else {
				if err := conn.WriteJSON(Message{FileContent: content, FilePath: msg.FilePath}); err != nil {
					log.Printf("[ERROR]发送文件内容失败: %v", err)
				}
			}

		case "writeFile":
			// 写入文件内容
			fileManager := GetFileManager()
			err := fileManager.WriteFile(msg.FilePath, msg.FileContent)
			if err != nil {
				sendError(conn, "[ERROR]写入文件失败: "+err.Error())
			} else {
				if err := conn.WriteJSON(Message{Status: "success", Command: "writeFile", FilePath: msg.FilePath}); err != nil {
					log.Printf("[ERROR]发送写入成功消息失败: %v", err)
				}
			}

		case "createDirectory":
			// 创建目录
			fileManager := GetFileManager()
			err := fileManager.CreateDirectory(msg.Content)
			if err != nil {
				sendError(conn, "[ERROR]创建目录失败: "+err.Error())
			} else {
				if err := conn.WriteJSON(Message{Status: "success", Command: "createDirectory"}); err != nil {
					log.Printf("[ERROR]发送创建成功消息失败: %v", err)
				}
			}

		case "deleteFile":
			// 删除文件或目录
			fileManager := GetFileManager()
			err := fileManager.DeleteFile(msg.Content)
			if err != nil {
				sendError(conn, "[ERROR]删除文件失败: "+err.Error())
			} else {
				if err := conn.WriteJSON(Message{Status: "success", Command: "deleteFile"}); err != nil {
					log.Printf("[ERROR]发送删除成功消息失败: %v", err)
				}
			}

		case "renameFile":
			// 重命名文件或目录
			fileManager := GetFileManager()
			err := fileManager.RenameFile(msg.OldPath, msg.NewPath)
			if err != nil {
				sendError(conn, "[ERROR]重命名文件失败: "+err.Error())
			} else {
				if err := conn.WriteJSON(Message{Status: "success", Command: "renameFile"}); err != nil {
					log.Printf("[ERROR]发送重命名成功消息失败: %v", err)
				}
			}
		case "downloadFile":
			fileManager := GetFileManager()
			fullPath := msg.FilePath
			if fullPath == "" {
				sendError(conn, "[ERROR]未指定文件路径")
				continue
			}
			// 读取文件内容
			content, err := fileManager.ReadFileRaw(fullPath)
			if err != nil {
				sendError(conn, "[ERROR]下载文件失败: "+err.Error())
				continue
			}
			b64 := base64.StdEncoding.EncodeToString(content)
			if err := conn.WriteJSON(Message{FileContent: b64, FilePath: msg.FilePath, Status: "download", Output: "download", Content: msg.FilePath}); err != nil {
				log.Printf("[ERROR]发送下载文件内容失败: %v", err)
			}
		case "uploadFile":
			fileManager := GetFileManager()
			if msg.FilePath == "" || msg.FileContent == "" {
				sendError(conn, "[ERROR]未指定文件路径或内容")
				continue
			}
			data, err := base64.StdEncoding.DecodeString(msg.FileContent)
			if err != nil {
				sendError(conn, "[ERROR]文件内容解码失败: "+err.Error())
				continue
			}
			err = fileManager.WriteFileRaw(msg.FilePath, data)
			if err != nil {
				sendError(conn, "[ERROR]上传文件失败: "+err.Error())
			} else {
				if err := conn.WriteJSON(Message{Status: "success", Command: "uploadFile", FilePath: msg.FilePath}); err != nil {
					log.Printf("[ERROR]发送上传成功消息失败: %v", err)
				}
			}

		case "createZip":
			// 创建zip文件
			fileManager := GetFileManager()
			if len(msg.FilesToZip) == 0 {
				sendError(conn, "[ERROR]未指定要压缩的文件")
				continue
			}
			if msg.ZipFileName == "" {
				sendError(conn, "[ERROR]未指定压缩文件名")
				continue
			}
			err := fileManager.CreateZipFile(msg.FilesToZip, msg.ZipFileName)
			if err != nil {
				sendError(conn, "[ERROR]创建zip文件失败: "+err.Error())
			} else {
				if err := conn.WriteJSON(Message{Status: "success", Command: "createZip", FilePath: msg.ZipFileName}); err != nil {
					log.Printf("[ERROR]发送压缩成功消息失败: %v", err)
				}
			}
		case "getConsoleHistory":
			// 处理终端历史获取命令
			systemMonitor := GetSystemMonitor()
			if systemMonitor != nil {
				history := systemMonitor.GetConsoleHistory()
				historyStr := ""
				if len(history) > 0 {
					historyStr = strings.Join(history, "") // 保持原有换行
				}
				if err := conn.WriteJSON(Message{FileContent: historyStr, Command: "getConsoleHistory"}); err != nil {
					log.Printf("[ERROR]发送终端历史失败: %v", err)
				}
			}
		case "getPanelConfig":
			// 读取Panel_Setting/config.json
			fileManager := GetFileManager()
			content, err := fileManager.ReadFile("Panel_Setting/config.json")
			if err != nil {
				sendError(conn, "读取配置失败: "+err.Error())
				break
			}
			if err := conn.WriteJSON(Message{FileContent: content}); err != nil {
				log.Printf("[ERROR]发送配置内容失败: %v", err)
			}
		case "setPanelConfig":
			// 写入Panel_Setting/config.json，msg.Content为json字符串
			fileManager := GetFileManager()
			if err := fileManager.WriteFile("Panel_Setting/config.json", msg.Content); err != nil {
				sendError(conn, "写入配置失败: "+err.Error())
				break
			}
			if err := conn.WriteJSON(Message{Status: "ok"}); err != nil {
				log.Printf("[ERROR]发送写入成功消息失败: %v", err)
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

	// 初始化文件管理器为ServerPath的父目录
	serverDir := config.ServerPath
	if serverDir != "" {
		parentDir := serverDir
		if stat, err := os.Stat(serverDir); err == nil && !stat.IsDir() {
			parentDir = filepath.Dir(serverDir)
		} else if err == nil && stat.IsDir() {
			parentDir = serverDir
		} else {
			log.Printf("[WARN] ServerPath无效，文件管理器根目录使用默认值")
		}
		InitFileManager(parentDir)
	} else {
		InitFileManager(".")
	}

	// 注册 WebSocket 路由，当访问 /ws 路径时，调用 handleWebSocket 函数处理
	http.HandleFunc("/ws", handleWebSocket)

	// 启动 HTTP 服务，监听配置文件中指定的端口
	log.Println("[ZephyCraft-Panel-2]websocket服务于端口" + config.Port + "上启动")
	// 启动服务，若失败则记录致命错误并终止程序（因 log.Fatal 会处理错误并退出）
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
